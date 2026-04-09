// Created by Sean L. on Apr. 8.
// Last Updated by Sean L. on Apr. 8.
//
// curio-api
// routes/submissions.go
//
// Makabaka1880, 2026. All rights reserved.

package routes

import (
	"curio-api/persistence"
	"curio-api/utils"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func generateToken() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func HandleSubmissionsRoute(r *gin.Engine) *gin.RouterGroup {
	submissions := r.Group("/submissions")
	submissions.POST("/upload", HandleSingleSubmission)
	submissions.DELETE("/retract/:id", HandleSingleDeletion)
	return submissions
}

func HandleSingleSubmission(c *gin.Context) {
	// MARK: Form Parsing
	name := c.PostForm("name")
	description := c.PostForm("description")
	event := c.PostForm("event")
	artifact, err := c.FormFile("artifact")
	if err != nil {
		c.JSON(400, gin.H{"error": "Artifact is required"})
		return
	}

	// Temp saving, ready for S3 upload
	objName := uuid.New()
	parts := strings.Split(artifact.Filename, ".")
	ext := parts[len(parts)-1]
	fileName := objName.String() + "." + ext
	dst := filepath.Join("./tmp/", fileName)
	contentType := "application/octet-stream"

	// MARK: S3 Upload
	c.SaveUploadedFile(artifact, dst)
	fmt.Printf("%v, %v", name, description)

	info, err := persistence.S3.FPutObject(utils.CTX, viper.GetString(utils.S3_BUCKET), fileName, dst, minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to upload artifact", "spec": err.Error()})
		return
	}
	_ = os.Remove(dst)

	eventID, err := uuid.Parse(event)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid event ID format", "spec": err.Error()})
		return
	}

	token := generateToken()

	persistence.DB.Create(&persistence.Submission{
		ID:             objName,
		EventID:        eventID,
		Name:           name,
		Description:    description,
		SubmissionDate: time.Now(),
		Token:          token,
	})

	c.JSON(200, gin.H{
		"message": "Submission successful",
		"data": gin.H{
			"id":          objName,
			"name":        name,
			"description": description,
			"token":       token,
			"artifactURL": fmt.Sprintf("https://%s/%s/%s", viper.GetString(utils.S3_ENDPOINT), viper.GetString(utils.S3_BUCKET), objName),
			"size":        info.Size,
		},
	})
}

func HandleSingleDeletion(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format", "spec": err.Error()})
		return
	}
	entry := &persistence.Submission{
		ID: parsedID,
	}
	err = persistence.DB.First(entry, "id=?", parsedID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Entry not found"})
			return
		} else {
			c.JSON(500, gin.H{"error": "Database error", "spec": err.Error()})
			return
		}
	}

	token := c.Query("token")
	if token == "" || token != entry.Token {
		c.JSON(400, gin.H{"error": "Invalid or missing token"})
		return
	}

	persistence.DB.Delete(&persistence.Submission{}, "id=?", id)
	persistence.S3.RemoveObject(utils.CTX, viper.GetString(utils.S3_BUCKET), entry.ID.String(),
		minio.RemoveObjectOptions{
			ForceDelete: true,
		},
	)
	c.JSON(200, gin.H{
		"message": "Entry deletion successful",
		"data": gin.H{
			"name": entry.Name,
		},
	})
}

func HandleEventSubmissionQueries(c *gin.Context) {

}
