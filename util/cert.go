package util

import (
	"io/ioutil"
	"log"
	"crypto/x509"
	"crypto/tls"
	"net/http"
)

// ReadCertKey - Read file csr and key
func ReadCertKey(pathCsr string, pathKey string) (caCert []byte, keyCert []byte){

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

func PreparTlsConfig(certFile string, keyFile string, openedCert []byte) (transporte *http.Transport){

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[INFO] - LoadedX509 cert: ", cert)

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(openedCert)


	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		InsecureSkipVerify: true,
	}

	log.Println("[INFO] - Created TlsConfig: ", tlsConfig)

	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return transport
}

