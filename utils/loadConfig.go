// Created by Sean L. on Apr. 8.
// Last Updated by Sean L. on Apr. 8.
//
// curio-api
// utils/loadConfig.go
//
// Makabaka1880, 2026. All rights reserved.

package utils

import "github.com/spf13/viper"

var (
	POSTGRES_USER     ViperSlot = "POSTGRES_USER"
	POSTGRES_PASSWORD ViperSlot = "POSTGRES_PASSWORD"
	POSTGRES_DB       ViperSlot = "POSTGRES_DB"
	POSTGRES_PORT     ViperSlot = "POSTGRES_PORT"
	POSTGRES_HOST     ViperSlot = "POSTGRES_HOST"
	S3_ENDPOINT       ViperSlot = "MINIO_ENDPOINT"
	S3_USER           ViperSlot = "MINIO_ROOT_USER"
	S3_PASSWORD       ViperSlot = "MINIO_ROOT_PASSWORD"
	S3_BUCKET         ViperSlot = "MINIO_BUCKET_NAME"
	AUTH_TOKEN        ViperSlot = "ADMIN_TOKEN"
)

type ViperSlot = string

func LoadConfig() {
	viper.AutomaticEnv()
}
