// Created by Sean L. on Apr. 8.
// Last Updated by Sean L. on Apr. 8.
//
// curio-api
// persistence/conn.go
//
// Makabaka1880, 2026. All rights reserved.

package persistence

import (
	"curio-api/utils"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var S3 *minio.Client

func ConnectDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString(utils.POSTGRES_HOST),
		viper.GetString(utils.POSTGRES_PORT),
		viper.GetString(utils.POSTGRES_USER),
		viper.GetString(utils.POSTGRES_PASSWORD),
		viper.GetString(utils.POSTGRES_DB),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Event{}, &Submission{})
	if err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return nil, err
	}

	log.Println("Database migration completed successfully.")

	DB = db
	return db, nil
}

func ConnectS3() (*minio.Client, error) {
	endpoint := viper.GetString(utils.S3_ENDPOINT)
	user := viper.GetString(utils.S3_USER)
	password := viper.GetString(utils.S3_PASSWORD)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(user, password, ""),
		Secure: false,
	})

	if err != nil {
		log.Printf("Failed to connect to S3: %v", err)
		return nil, err
	} else {
		S3 = minioClient
	}
	return S3, nil
}
