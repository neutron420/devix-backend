package media

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers media routes (static file serving).
func RegisterRoutes(r *gin.Engine, uploadDir string) {
	// Serve uploaded files at /uploads/*
	r.Static("/uploads", uploadDir)
}
