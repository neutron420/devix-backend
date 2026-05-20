package org

import (
	"devix-backend/internal/middleware"
	jwtpkg "devix-backend/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler, jwtManager *jwtpkg.Manager) {
	orgs := rg.Group("/organizations")
	{
		orgs.GET("", handler.List)
		orgs.POST("", middleware.Auth(jwtManager), handler.CreateOrg)
		orgs.GET("/slug/:slug", handler.GetBySlug)
		orgs.GET("/:id", handler.GetByID)
		orgs.GET("/:id/members", handler.GetMembers)
		orgs.POST("/:id/members", middleware.Auth(jwtManager), handler.AddMember)
		orgs.DELETE("/:id/members/:user_id", middleware.Auth(jwtManager), handler.RemoveMember)
		orgs.PUT("/:id", middleware.Auth(jwtManager), handler.UpdateOrg)
		orgs.DELETE("/:id", middleware.Auth(jwtManager), handler.DeleteOrg)
	}
}
