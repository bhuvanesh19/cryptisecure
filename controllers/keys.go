package controllers

import (
	"cryptisecure/models"
	"cryptisecure/validators"

	"github.com/gin-gonic/gin"
)

func AddKey(ctx *gin.Context) {
	var req validators.AddKey
	ctx.Bind(&req)
	if models.KeyExists(req.Id) {
		logger.Println("Key already exists")
		ctx.JSON(400, map[string]interface{}{
			"error": "Key already exists",
		})
		return
	}
	key_data, err := req.Key.Open()
	if err != nil {
		logger.Println(err)
		ctx.JSON(400, map[string]interface{}{
			"error": "Invalid key file",
		})
		return
	}
	var key_bytes []byte = make([]byte, req.Key.Size)
	_, err = key_data.Read(key_bytes)
	if err != nil {
		logger.Println(err)
		ctx.JSON(400, map[string]interface{}{
			"error": "Failed to read key file",
		})
		return
	}
	certificate, err := models.GetCertificate(req.Id)
	if err != nil {
		logger.Println(err)
		ctx.JSON(404, map[string]interface{}{
			"error": "Certificate not found",
		})
		return
	}

	err = models.AddPrivateKey(certificate, key_bytes)
	if err != nil {
		logger.Println(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to add private key",
		})
		return
	}

}

func DeleteKey(ctx *gin.Context) {
	id := ctx.Params.ByName("id")

	err := models.DeletePrivateKey(id)
	if err != nil {
		logger.Println(err)
		ctx.JSON(500, map[string]interface{}{
			"error": "Failed to delete key",
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"message": "Key deleted successfully",
	})
}
