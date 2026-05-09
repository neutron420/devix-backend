package tag

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	tags := rg.Group("/tags")
	{
		tags.GET("", handler.GetAll)
		tags.GET("/trending", handler.GetTrending)
	}
}
