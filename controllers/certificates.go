package controllers

import (
	"cryptisecure/models"
	"cryptisecure/validators"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var logger = log.Default()

func GetCertificate(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	cert, err := models.GetCertificate(id)
	if err != nil {
		ctx.JSON(404, map[string]interface{}{
			"error": "Certificate not found",
		})
		return
	}
	ctx.JSON(200, cert)
}

func ListCertificate(ctx *gin.Context) {
	var req validators.GetCertificateRequest
	if err := ctx.BindQuery(&req); err != nil {
		logger.Println(err)
		ctx.JSON(400, map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}
	certificates, err := models.ListCertificates(req.ID, req.Name, req.Group, req.Domin)
	if err != nil {
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to retrieve certificates",
		})
		return
	}
	ctx.JSON(200, certificates)

}

func AddCertificate(ctx *gin.Context) {
	var req validators.AddCertificateRequest
	ctx.Bind(&req)

	file, err := req.File.Open()
	if err != nil {
		logger.Fatal(err)
		ctx.JSON(400, map[string]interface{}{
			"error": "Failed to read file",
		})
		return
	}

	file_content := make([]byte, req.File.Size)
	file.Read(file_content)
	pem_decoded_block, _ := pem.Decode(file_content)
	if pem_decoded_block == nil || pem_decoded_block.Type != "CERTIFICATE" {
		ctx.JSON(400, map[string]interface{}{
			"error": "Invalid certificate file. Allowed formats : pem",
		})
		return
	}
	cert, err := x509.ParseCertificate(pem_decoded_block.Bytes)
	if err != nil {
		ctx.JSON(400, map[string]interface{}{
			"error": "Invalid certificate. Make sure to upload a valid X509 certificate",
		})
		return
	}

	var certificate models.Certificate = models.Certificate{
		ID:     cert.SerialNumber.String(),
		Name:   req.Name,
		Group:  req.Group,
		Domain: req.Domain,
	}

	if cert_in_db, err := models.IDIntegrityCheck(certificate.ID); err != nil || !cert_in_db {
		logger.Print("Cert in DB ", cert_in_db)
		logger.Print(err)
		ctx.JSON(400, map[string]interface{}{
			"error": "Integrity Validation Failed",
		})
		return
	}

	if err := models.AddCertificate(&certificate); err != nil {
		logger.Fatal(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to add certificate",
		})
		return
	}
	os.WriteFile(fmt.Sprintf("certificates/%s.%s", certificate.ID, "pem"), file_content, os.ModePerm.Perm())
	ctx.JSON(200, certificate)
}

func DeleteCertificate(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	cert, err := models.GetCertificate(id)
	if err != nil {
		ctx.JSON(404, map[string]interface{}{
			"error": "Certificate not found",
		})
		return
	}
	if err := os.Remove(fmt.Sprintf("certificates/%s.%s", cert.ID, "pem")); err != nil {
		logger.Println(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to delete certificate",
		})
		return
	}
	if err := models.DeleteCertificate(id); err != nil {
		logger.Fatal(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to delete certificate from database",
		})
		return
	}
	ctx.JSON(200, map[string]interface{}{
		"message": "Certificate deleted successfully",
	})
}
