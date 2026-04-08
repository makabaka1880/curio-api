// Created by Sean L. on Apr. 8.
// Last Updated by Sean L. on Apr. 8.
//
// curio-api
// middleware/auth.go
//
// Makabaka1880, 2026. All rights reserved.

package middleware

import (
	"curio-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func AuthMiddleware(c *gin.Context) {
	if c.GetHeader("Authentication") != viper.GetString(utils.AUTH_TOKEN) {
		c.JSON(403, gin.H{"error": "Forbidden"})
		c.Abort()
	}
}
