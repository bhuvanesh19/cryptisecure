package models

import (
	"log"
)

type Certificate struct {
	ID     string `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Group  string `json:"group"`
	Domain string `json:"domain"`
}

func ListCertificates(id, name, group, domain string) ([]Certificate, error) {
	log.Println("Query str", id, name, group, domain)
	var query = DB.Model(&Certificate{})
	if id != "" {
		log.Println("Id used for querying")
		query = query.Where("id = ?", id)
	}
	if name != "" {
		log.Println("Name used for querying")
		query = query.Where("name = ?", name)
	}
	if group != "" {
		log.Println("Group used for querying")
		query = query.Where("group = ?", group)
	}
	if domain != "" {
		log.Println("Domain used for querying")
		query = query.Where("domain = ?", domain)
	}

	var certificates []Certificate
	err := query.Find(&certificates).Error
	log.Print("SQl", query.Statement.SQL.String())
	return certificates, err
}

func GetCertificate(id string) (*Certificate, error) {
	var certificate Certificate
	err := DB.Where("id =?", id).First(&certificate).Error
	log.Println("Search for certificate with id", id, certificate)
	return &certificate, err
}

func IDIntegrityCheck(id string) (bool, error) {
	var count int64
	err := DB.Model(&Certificate{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func AddCertificate(certificate *Certificate) error {
	return DB.Create(certificate).Error
}

func DeleteCertificate(id string) error {
	return DB.Where("id = ?", id).Delete(&Certificate{}).Error
}
