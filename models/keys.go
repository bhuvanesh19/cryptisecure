package models

import (
	"crypto/tls"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
)

type CerKey struct {
	ID  string `json:"id" gorm:"primaryKey"`
	Key string `json:"key"`
}

func AddPrivateKey(certificate *Certificate, keyPEM []byte) error {
	block, _ := pem.Decode(keyPEM)
	if block == nil || block.Type != "PRIVATE KEY" {
		return errors.New("failed to decode private key PEM")
	}

	cert_file, err := os.ReadFile(fmt.Sprintf("certificates/%s.%s", certificate.ID, "pem"))
	if err != nil {
		return fmt.Errorf("failed to read certificate file: %w", err)
	}

	_, err = tls.X509KeyPair(cert_file, []byte(keyPEM))
	if err != nil {
		return fmt.Errorf("certificate and private key do not match: %w", err)
	}

	if err := DB.Save(&CerKey{
		certificate.ID,
		string(keyPEM),
	}).Error; err != nil {
		log.Println("Error saving private key ", err)
		return err
	}
	return nil
}

func DeletePrivateKey(id string) error {
	return DB.Where("id =?", id).Delete(&CerKey{}).Error
}

func GetKey(id string) (*CerKey, error) {
	var cerKey CerKey
	err := DB.Where("id =?", id).First(&cerKey).Error
	return &cerKey, err
}

func KeyExists(id string) bool {
	var cerKey CerKey
	err := DB.Where("id =?", id).First(&cerKey).Error
	return err == nil
}
