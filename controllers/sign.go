package controllers

import (
	"archive/tar"
	"bytes"
	"cryptisecure/models"
	"cryptisecure/utils"
	"cryptisecure/validators"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"log"

	"github.com/gin-gonic/gin"
)

func SignFile(ctx *gin.Context) {
	var req validators.SignFileRequest
	if err := ctx.Bind(&req); err != nil {
		logger.Println(err)
		ctx.JSON(400, map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	// Read the file
	file, err := req.File.Open()
	if err != nil {
		logger.Println(err)
		ctx.JSON(400, map[string]interface{}{
			"error": "Failed to open file",
		})
		return
	}
	var file_bytes []byte = make([]byte, req.File.Size)
	_, err = file.Read(file_bytes)
	if err != nil {
		logger.Println(err)
		ctx.JSON(400, map[string]interface{}{
			"error": "Failed to read file",
		})
		return
	}

	// Read the key
	key, err := models.GetKey(req.ID)
	if err != nil {
		logger.Println(err)
		ctx.JSON(404, map[string]interface{}{
			"error": "Key not found",
		})
		return
	}
	key_bytes := []byte(key.Key)
	key_block, _ := pem.Decode(key_bytes)
	if key_block == nil || key_block.Type != "PRIVATE KEY" {
		logger.Println("Invalid key format")
		ctx.JSON(400, map[string]interface{}{
			"error": "Invalid key format",
		})
		return
	}

	// Write the file
	var tar_file_buffer bytes.Buffer
	tar_writer := tar.NewWriter(&tar_file_buffer)

	err = tar_writer.WriteHeader(&tar.Header{
		Name: req.File.Filename,
		Size: int64(len(file_bytes)),
		Mode: 0644,
	})
	if err != nil {
		logger.Println(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to write to tar file",
		})
		return
	}
	_, err = tar_writer.Write(file_bytes)
	if err != nil {
		logger.Println(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to write to tar file",
		})
		return
	}

	// prepare the metadata json file. Write the ID
	var metadata_json map[string]interface{} = make(map[string]interface{})
	metadata_json["id"] = req.ID

	// Calculate the signature
	file_hash := utils.CalculateSHA256(file_bytes)
	log.Println("Hash Calculated: ", file_hash)
	privateKey, err := x509.ParsePKCS8PrivateKey(key_block.Bytes)
	if err != nil {
		logger.Println(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to parse private key",
		})
		return
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		logger.Println("Invalid private key format")
		ctx.JSON(400, map[string]interface{}{
			"error": "Invalid private key format",
		})
		return
	}

	// Encrypt the data using the private key
	// Note: Use a proper padding scheme (e.g., PKCS1v15)
	encryptedData, err := rsa.SignPKCS1v15(
		rand.Reader,
		rsaPrivateKey,
		crypto.SHA256,
		file_hash,
	)
	log.Println("Encrypted data: ", encryptedData)
	if err != nil {
		logger.Println(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to sign data",
		})
		return
	}

	metadata_json["signature"] = base64.StdEncoding.EncodeToString(encryptedData)

	// Write the signature to the metadata json file
	metadata_json_bytes, err := json.MarshalIndent(metadata_json, "", "  ")
	if err != nil {
		logger.Println(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to marshal metadata json",
		})
		return
	}
	err = tar_writer.WriteHeader(&tar.Header{
		Name: "metadata.json",
		Size: int64(len(metadata_json_bytes)),
		Mode: 0644,
	})
	if err != nil {
		logger.Println(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to write to tar file",
		})
		return
	}
	_, err = tar_writer.Write(metadata_json_bytes)
	if err != nil {
		logger.Println(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to write to tar file",
		})
		return
	}
	err = tar_writer.Flush()
	if err != nil {
		logger.Println(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to flush tar writer",
		})
		return
	}

	tar_writer.Close()
	ctx.Header("Content-Type", "application/x-tar")
	ctx.Header("Content-Disposition", "attachment; filename=archive.tar")

	// Write the tar file to the response
	ctx.Data(200, "application/x-tar", tar_file_buffer.Bytes())
}
