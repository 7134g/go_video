package proxy

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"

	"github.com/google/martian/mitm"
)

func (p *Proxy) loadCert() error {
	ca, err := loadRootCA(p.config.CAFile)
	if err != nil {
		return err
	}
	p.ca = ca

	key, err := loadRootKey(p.config.PrivateKeyFile)
	if err != nil {
		return err
	}
	p.privateKey = key

	return nil
}

func loadRootCA(filename string) (*x509.Certificate, error) {
	mitmData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(mitmData)
	return x509.ParseCertificate(block.Bytes)
}

func loadRootKey(filename string) (*rsa.PrivateKey, error) {
	priData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(priData)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func GenMITM(caFile, keyFile string) error {
	caData, priData, err := newCert(DefaultDomain, DefaultOrganization, DefaultValidDuration)
	if err != nil {
		return err
	}

	if err := os.WriteFile(keyFile, priData, 0600); err != nil {
		return err
	}

	if err := os.WriteFile(caFile, caData, 0644); err != nil {
		return err
	}

	return nil
}

func newCert(domain string, organization string, validDuration time.Duration) ([]byte, []byte, error) {
	pub, prv, err := mitm.NewAuthority(domain, organization, validDuration)
	if err != nil {
		return nil, nil, err
	}
	prvData := bytes.NewBuffer([]byte{})
	if err = pem.Encode(prvData, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(prv)}); err != nil {
		return nil, nil, err
	}
	certData := bytes.NewBuffer([]byte{})
	if err = pem.Encode(certData, &pem.Block{Type: "CERTIFICATE", Bytes: pub.Raw}); err != nil {
		return nil, nil, err
	}

	return certData.Bytes(), prvData.Bytes(), nil
}
