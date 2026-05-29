package proxy

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
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

// CheckCertInstalled 探测当前 CA 是否已加入系统信任存储。
// Windows: certutil；macOS: security；Linux: 扫描常见 ca-bundle 比对 DER 指纹。
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
	case "linux":
		return checkLinuxBundle(ca), nil
	default:
		return false, nil
	}
}

// checkLinuxBundle 在常见发行版的 CA bundle 里按 SHA-256 指纹查找当前 CA。
// 系统层信任就足以判定 installer 跑过——NSS 是浏览器层增强,失败仍能从这里看到日志。
func checkLinuxBundle(ca *x509.Certificate) bool {
	target := sha256.Sum256(ca.Raw)
	for _, path := range []string{
		"/etc/ssl/certs/ca-certificates.crt", // Debian/Ubuntu/Alpine
		"/etc/pki/tls/certs/ca-bundle.crt",   // RHEL/Fedora/CentOS
		"/etc/ssl/cert.pem",                  // Arch
	} {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		rest := data
		for {
			var block *pem.Block
			block, rest = pem.Decode(rest)
			if block == nil {
				break
			}
			if block.Type != "CERTIFICATE" {
				continue
			}
			if sha256.Sum256(block.Bytes) == target {
				return true
			}
		}
	}
	return false
}
