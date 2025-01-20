package main

import (
	"cryptisecure/models"
	"cryptisecure/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	models.Setup()

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run(":8080")

}
