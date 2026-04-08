// Created by Sean L. on Apr. 8.
// Last Updated by Sean L. on Apr. 8.
//
// curio-api
// utils/context.go
//
// Makabaka1880, 2026. All rights reserved.

package utils

import "context"

var CTX context.Context

func init() {
	CTX = context.Background()
}
