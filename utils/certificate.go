package utils

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func ReadCertificate(id string) (*x509.Certificate, error) {
	certData, err := os.ReadFile(fmt.Sprintf("certificates/%s.%s", id, "pem"))
	if err != nil {
		log.Println("failed to read certificate file: ", err)
		return nil, err
	}

	// Decode the certificate from PEM format
	block, _ := pem.Decode(certData)
	if block == nil {
		log.Println("failed to decode PEM block from certificate")
		return nil, err

	}

	// Parse the X.509 certificate
	return x509.ParseCertificate(block.Bytes)
}
