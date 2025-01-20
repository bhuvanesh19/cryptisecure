package routes

import (
	"cryptisecure/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/certificate", controllers.ListCertificate)
	r.GET("/certificate/:id", controllers.GetCertificate)
	r.POST("/certificate", controllers.AddCertificate)
	r.DELETE("/certificate/:id", controllers.DeleteCertificate)

	r.POST("/key", controllers.AddKey)
	r.DELETE("/key/:id", controllers.DeleteKey)

	r.POST("/sign", controllers.SignFile)
	r.POST("/verify", controllers.VerifyFile)
}
