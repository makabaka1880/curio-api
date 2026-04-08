package main

import (
	"curio-api/persistence"
	"curio-api/routes"
	"curio-api/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	utils.LoadConfig()
	persistence.ConnectDB()
	persistence.ConnectS3()
	routes.HandleSubmissionsRoute(r)
	routes.HandleEventsRoutes(r)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
