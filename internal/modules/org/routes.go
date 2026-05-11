package org

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtManager *jwtpkg.Manager) {
	orgs := rg.Group("/organizations")
	{
		orgs.POST("", middleware.Auth(jwtManager), handler.CreateOrg)
		orgs.GET("/:id/members", handler.GetMembers)
		orgs.POST("/:id/members", middleware.Auth(jwtManager), handler.AddMember)
	}
}
