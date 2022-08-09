// certificate is a CA root certificate factory

package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sequencer"
	"time"
)

const (
	RSAKeySize = 4096
)

// CertificateConfig provides required parameters for a CA certificate
type CertificateConfig struct {
	Country      string   `json:"country"`
	CommonName   string   `json:"common_name"`
	Organization []string `json:"organisation"`
	AltNames     AltNames `json:"alt_names"`
}

type AltNames struct {
	DNSNames []string `json:"dns_names"`
	IPs      []net.IP `json:"ips"`
}

// Certificate stores CA certificate wih input parameters
type Certificate struct {
	Host           string    `json:"host"`
	ValidFor       int       `json:"valid_for"`
	Unit           string    `json:"unit"`
	ValidSince     time.Time `json:"valid_since"`
	ValidTill      time.Time `json:"valid_till"`
	CertificatePEM []byte    `json:"certificate_pem"`
	PrivateKeyPEM  []byte    `json:"private_key_pem"`
	PublicKeyPREM  []byte    `json:"public_key_pem"`
}

// NewPrivateKey makes a new private key
func NewPrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, RSAKeySize)
}

// EncodePublicKeyPEM serializes public key to a seq of bytes
func EncodePublicKeyAsPEM(key *rsa.PublicKey) ([]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return []byte{}, err
	}
	block := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}
	return pem.EncodeToMemory(&block), nil
}

// EncodePrivateKeyAsPEM serializes private key to a seq of bytes
func EncodePrivateKeyAsPEM(key *rsa.PrivateKey) []byte {
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return pem.EncodeToMemory(&block)
}

// EncodeCertificateAsPEM serializes certificate to a seq of bytes
func EncodeCertificateAsPEM(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}

// NewSignedCertificate fills up the x509 structure with prepared values
func NewRootCertificate(cfg CertificateConfig, host string, seqNo int,
	validSince, validTille time.Time,
	publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) (*x509.Certificate, error) {
	tmplate := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(seqNo),
		Subject: pkix.Name{
			Country:      cfg.Country,
			CommonName:   cfg.CommonName,
			Organization: cfg.Organization,
		},
		NotBefore:             validSince,
		NotAfter:              validTill,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{host},
	}

	certDERBytes, err := x509.CreateCertificate(rand.Reader, &tmplate, &tmplate, publicKey, privateKey)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(certDERBytes)
}

// NewCertificate creates new certificate using config for a host for a time with a seq no
func NewCertificate(cfg CertificateConfig, host string, validFor int, unit string) (*Certificate, error) {
	seqNo, err := sequencer.NextSeqNo()
	if err != nil {
		return nil, err
	}

	validSince, validTill, err := timerange.GetValidityPriod(validFor, unit)
	if err != nil {
		return nil, err
	}

	privateKey, publicKey, err := generateKeyPair(RSAKeySize)
	if err != nil {
		return nil, err
	}

	cert, err := NewRootCertificate(cfg, host, seqNo, validSince, validTill, publicKey, privateKey)
	if err != nil {
		return nil, err
	}

	cert := &Certificate{
		Host:           host,
		ValidFor:       validFor,
		Unit:           unit,
		ValidSince:     validSince,
		ValidTill:      validTill,
		CertificatePEM: EncodeCertificatePEM(cert),
		PrivateKeyPEM:  EncodePrivateKeyPEM(privateKey),
		PublicKeyPEM:   EncodePublicKeyAsPEM(publicKey),
	}

	return cert, nil
}

// Certificate is a Stringer meaning it can be printed as Json
func (cert *Certificate) StringAsJson() string {
	out, err := json.Marshal(cert)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(string(out))
}
