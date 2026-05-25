package post

import (
	"context"
	"errors"
	"strings"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/middleware"
	"devix-backend/internal/modules/analytics"
	"devix-backend/internal/modules/media"
	"devix-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service      *Service
	mediaService *media.Service
	analyticsService *analytics.Service
}

func NewHandler(service *Service, mediaService *media.Service, analyticsService *analytics.Service) *Handler {
	return &Handler{service: service, mediaService: mediaService, analyticsService: analyticsService}
}

func (h *Handler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid request: "+err.Error()))
		return
	}
	result, err := h.service.Create(c.Request.Context(), userID, &req)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.Created(c, result)
}

func (h *Handler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	result, err := h.service.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OK(c, result)

	if h.analyticsService == nil {
		return
	}
	postID, err := uuid.Parse(result.ID)
	if err != nil {
		return
	}
	ua := c.GetHeader("User-Agent")
	referrer := c.GetHeader("Referer")
	if referrer == "" {
		referrer = c.GetHeader("Referrer")
	}
	country := strings.TrimSpace(c.GetHeader("CF-IPCountry"))
	if country == "" {
		country = strings.TrimSpace(c.GetHeader("X-Country"))
	}
	ip := c.ClientIP()
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		h.analyticsService.TrackView(ctx, postID, ua, ip, country, referrer)
	}()
}

func (h *Handler) ListFollowing(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}
	var query FeedQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid query parameters"))
		return
	}
	if query.After == "" && query.Before == "" && query.Cursor != "" {
		query.After = query.Cursor
	}
	if query.After != "" && query.Before != "" {
		response.Error(c, apperrors.BadRequest("Use only one of 'after' or 'before'"))
		return
	}
	result, err := h.service.ListFollowing(c.Request.Context(), userID, query)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OKWithMeta(c, result.Posts, &response.Meta{Cursor: result.Cursor, PrevCursor: result.PrevCursor, HasMore: result.HasMore})
}

func (h *Handler) List(c *gin.Context) {
	var query FeedQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid query parameters"))
		return
	}
	if query.After == "" && query.Before == "" && query.Cursor != "" {
		query.After = query.Cursor
	}
	if query.After != "" && query.Before != "" {
		response.Error(c, apperrors.BadRequest("Use only one of 'after' or 'before'"))
		return
	}

	// Set RequestUserID for personalized ranking if logged in
	if userID, ok := middleware.GetUserID(c); ok {
		query.RequestUserID = userID
	}

	result, err := h.service.List(c.Request.Context(), query)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OKWithMeta(c, result.Posts, &response.Meta{Cursor: result.Cursor, PrevCursor: result.PrevCursor, HasMore: result.HasMore})
}

func (h *Handler) ListExplore(c *gin.Context) {
	var query FeedQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid query parameters"))
		return
	}
	if query.After == "" && query.Before == "" && query.Cursor != "" {
		query.After = query.Cursor
	}
	if query.After != "" && query.Before != "" {
		response.Error(c, apperrors.BadRequest("Use only one of 'after' or 'before'"))
		return
	}

	var userID uuid.UUID
	if uid, ok := middleware.GetUserID(c); ok {
		userID = uid
	}

	result, err := h.service.ListExplore(c.Request.Context(), userID, query)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OKWithMeta(c, result.Posts, &response.Meta{Cursor: result.Cursor, PrevCursor: result.PrevCursor, HasMore: result.HasMore})
}


func (h *Handler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}
	slug := c.Param("slug")
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid request: "+err.Error()))
		return
	}
	result, err := h.service.Update(c.Request.Context(), slug, userID, &req)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OK(c, result)
}

func (h *Handler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}
	slug := c.Param("slug")
	role := c.GetString(middleware.ContextKeyUserRole)
	if err := h.service.Delete(c.Request.Context(), slug, userID, role); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.NoContent(c)
}

func (h *Handler) UploadMedia(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}
	slug := c.Param("slug")

	post, err := h.service.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	if post.Author == nil || post.Author.ID != userID.String() {
		response.Error(c, apperrors.Forbidden("You can only add media to your own posts"))
		return
	}
	postID, err := uuid.Parse(post.ID)
	if err != nil {
		response.Error(c, apperrors.Internal(err))
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid form data"))
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		response.Error(c, apperrors.BadRequest("No files uploaded"))
		return
	}
	var uploads []*media.FileUpload
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			response.Error(c, apperrors.BadRequest("Failed to read file"))
			return
		}
		defer f.Close()
		uploads = append(uploads, &media.FileUpload{File: f, Header: fh})
	}
	uploaded, err := h.mediaService.UploadPostMedia(c.Request.Context(), postID, uploads)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.Created(c, uploaded)
}

func (h *Handler) ListDrafts(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}
	drafts, err := h.service.ListDrafts(c.Request.Context(), userID)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OK(c, drafts)
}

func (h *Handler) Autosave(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}
	slug := c.Param("slug")
	var req AutosaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid request: "+err.Error()))
		return
	}
	if err := h.service.Autosave(c.Request.Context(), slug, userID, &req); err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			response.Error(c, appErr)
			return
		}
		response.Error(c, apperrors.Internal(err))
		return
	}
	response.OK(c, gin.H{"message": "draft autosaved"})
}
