//go:generate go run update.go

package gobundledhttp

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"

	"github.com/tydavis/gobundledhttp/certificates"
)

var pool *x509.CertPool

func init() {
	// Always build the pool
	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(certificates.PemCerts) // from certificates.go file
}

// NewClient returns a new http.Client with the included certpool.
func NewClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
				RootCAs:    pool,
			},
			DisableCompression: true,
		},
	}
}

// MaxClient returns a maximum-security client with highly restrictive TLS settings
func MaxClient() *http.Client {
	config := tls.Config{
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
		RootCAs: pool,
	}

	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:    &config,
			DisableCompression: true,
		},
	}
}

// InsecureClient without certificate checking.
func InsecureClient() *http.Client {
	// Insecure client without cert-trust checking
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: true,
			},
			DisableCompression: true,
		},
	}

}

// GetPool returns the certificate pool.
func GetPool() *x509.CertPool {
	return pool
}

// AddPEM adds a PEM-formatted byte slice to certificate pool
func AddPEM(b []byte) bool {
	return pool.AppendCertsFromPEM(b)
}
