package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"os"
)

// CreateSSLCerts creates a new TLS configuration with the provided certificates
func NewCerts(certFile, keyFile, caFile string, insecure bool) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("failed to append CA certs")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	if insecure {
		tlsConfig.InsecureSkipVerify = true
	} else {
		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}
