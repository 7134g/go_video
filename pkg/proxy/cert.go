package proxy

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/google/martian/mitm"
)

const (
	DefaultDomain        = "proxy.local"
	DefaultOrganization  = "Proxy CA"
	DefaultValidDuration = 365 * 24 * time.Hour
	CACertFile           = "ca.crt"
	CAKeyFile            = "ca.key"
)

func GenCA() error {
	caData, priData, err := newCert(DefaultDomain, DefaultOrganization, DefaultValidDuration)
	if err != nil {
		return err
	}
	if err := os.WriteFile(CAKeyFile, priData, 0600); err != nil {
		return err
	}
	return os.WriteFile(CACertFile, caData, 0644)
}

func LoadCA() (*x509.Certificate, error) {
	data, err := os.ReadFile(CACertFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	return x509.ParseCertificate(block.Bytes)
}

func LoadKey() (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(CAKeyFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func newCert(domain, org string, validDuration time.Duration) ([]byte, []byte, error) {
	pub, prv, err := mitm.NewAuthority(domain, org, validDuration)
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

func CheckCertInstalled() (bool, error) {
	ca, err := LoadCA()
	if err != nil {
		return false, err
	}

	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("certutil", "-verifystore", "Root", ca.Subject.CommonName)
		return cmd.Run() == nil, nil
	case "darwin":
		cmd := exec.Command("security", "find-certificate", "-c", ca.Subject.CommonName, "-p", "/Library/Keychains/System.keychain")
		return cmd.Run() == nil, nil
	default:
		return false, nil
	}
}
