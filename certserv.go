package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"time"
)

func main() {
	// hostname for the certificate and the request
	const host = "foo.com"

	// create a new cert for 'foo.com' to be used in the HTTPS server
	certPEM, keyPEM := makeCert(host)
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Fatal("tls.X509KeyPair: ", err)
	}

	// start an HTTPS server with that cert
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tlsConfig.BuildNameToCertificate()

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	srv.Config.TLSConfig = tlsConfig
	srv.StartTLS()
	defer srv.Close()

	// create a client for the request which has
	// the cert as the only rootCA
	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM([]byte(certPEM))

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: rootCAs,
				// InsecureSkipVerify: true, // this will make it 'work' but then the hostname isn't checked
			},
		},
	}

	// make a request to the server
	req, err := http.NewRequest("GET", srv.URL, nil)
	if err != nil {
		log.Fatal("http.NewRequest: ", err)
	}
	req.Header.Set("Host", host)

	// this should work but always fails with
	// Get https://127.0.0.1:49805: x509: certificate signed by unknown authority (possibly because of "crypto/rsa: verification error" while trying to verify candidate authority certificate "serial:1000")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("client.Do: ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("ioutil.ReadAll: ", err)
	}

	// this should print '200 OK'
	fmt.Printf("%d %s\n", resp.StatusCode, string(body))
}

// makeCert creates a self-signed certificate for an HTTPS server
// which is valid for 'host' for one minute.
func makeCert(host string) (certPEM, keyPEM []byte) {
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatal("rsa.GenerateKey: ", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1000),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Minute),

		IsCA:                  true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{host},
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		log.Fatal("x509.CreateCertificate: ", err)
	}

	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert})
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	return
}
