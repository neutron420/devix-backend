package media

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, uploadDir string) {

	r.Static("/uploads", uploadDir)
}
