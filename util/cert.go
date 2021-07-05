package util

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/logger"
)

// ReadCertKey - Read file csr and key
func ReadCertKey(pathCsr string, pathKey string) (caCert []byte, keyCert []byte) {

	caCert, err := ioutil.ReadFile(pathCsr)
	if err != nil {
		log.Fatal(err)
	}

	keyCert, err = ioutil.ReadFile(pathKey)
	if err != nil {
		log.Fatal(err)
	}

	return caCert, keyCert
}

func PreparTlsConfig(certFile string, keyFile string, openedCert []byte) (transporte *http.Transport) {

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("LoadedX509 cert: ", cert)

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(openedCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

	logger.Info("Created TlsConfig: ", tlsConfig)

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return transport
}
