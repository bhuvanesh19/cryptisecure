package controllers

import (
	"archive/tar"
	"bytes"
	"cryptisecure/utils"
	"cryptisecure/validators"
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

func VerifyFile(ctx *gin.Context) {
	var req validators.VerifyFile
	ctx.Bind(&req)

	file, err := req.File.Open()
	if err != nil {
		logger.Println(err)
		ctx.JSON(400, map[string]interface{}{
			"error": "Failed to read file",
		})
		return
	}
	var file_data []byte = make([]byte, req.File.Size)
	_, err = file.Read(file_data)
	if err != nil {
		logger.Println(err)
		ctx.JSON(400, map[string]interface{}{
			"error": "Failed to read file",
		})
		return
	}

	//
	reader := bytes.NewReader(file_data)

	tar_reader := tar.NewReader(reader)
	var metadata_json map[string]interface{} = make(map[string]interface{})
	var file_hash []byte

	for i := 0; i < 2; i++ {
		header, err := tar_reader.Next()
		if err == io.EOF {
			log.Println("EOF Reached parsing Tar")
			break
		}
		if err != nil {
			logger.Println(err)
			ctx.JSON(500, map[string]interface{}{
				"error": "Failed to read tar file",
			})
			return
		}
		log.Println("Header ", header)
		if header.Name == "metadata.json" {
			var metadataBytes bytes.Buffer
			if _, err := io.Copy(&metadataBytes, tar_reader); err != nil {
				logger.Println(err)
				ctx.JSON(500, map[string]interface{}{
					"error": "Failed to read metadata.json",
				})
				return
			}
			err = json.Unmarshal(metadataBytes.Bytes(), &metadata_json)
			if err != nil {
				logger.Println(err)
				ctx.JSON(500, map[string]interface{}{
					"error": "Failed to unmarshal metadata.json",
				})
				return
			}
		} else {
			var file_data bytes.Buffer
			if _, err := io.Copy(&file_data, tar_reader); err != nil {
				logger.Println(err)
				ctx.JSON(500, map[string]interface{}{
					"error": fmt.Sprintf("Failed to read %s", header.Name),
				})
				return
			}
			file_hash = utils.CalculateSHA256(file_data.Bytes())
			log.Println("Hash calculated", file_hash)
		}
	}
	logger.Println("Metadata Json", metadata_json)
	id, ok := metadata_json["id"].(string)
	if !ok {
		logger.Println(err)
		ctx.JSON(400, map[string]interface{}{
			"error": "Invalid metadata.json",
		})
		return
	}

	signature, ok := metadata_json["signature"].(string)
	if !ok {
		logger.Println(err)
		ctx.JSON(400, map[string]interface{}{
			"error": "Invalid metadata.json",
		})
		return
	}
	decoded_signature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		logger.Println(err)
		ctx.JSON(400, map[string]interface{}{
			"error": "Invalid signature format",
		})
		return
	}

	certificate, err := utils.ReadCertificate(id)
	if err != nil {
		logger.Println(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to read certificate",
		})
		return
	}

	pubKey, ok := certificate.PublicKey.(*rsa.PublicKey)
	if !ok {
		logger.Println(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Invalid certificate format",
		})
		return
	}
	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, file_hash, decoded_signature)
	if err != nil {
		log.Println(err)
		ctx.JSON(400, map[string]interface{}{
			"error": "Invalid signature",
		})
		return
	}
	ctx.JSON(200, map[string]interface{}{
		"verified": true,
	})

}
