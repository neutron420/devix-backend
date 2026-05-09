package post

import (
	"errors"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/middleware"
	"devix-backend/internal/modules/media"
	"devix-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service      *Service
	mediaService *media.Service
}

func NewHandler(service *Service, mediaService *media.Service) *Handler {
	return &Handler{service: service, mediaService: mediaService}
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
	slug := c.Param("id")
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
	response.OKWithMeta(c, result.Posts, &response.Meta{Cursor: result.Cursor, HasMore: result.HasMore})
}

func (h *Handler) List(c *gin.Context) {
	var query FeedQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid query parameters"))
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
	response.OKWithMeta(c, result.Posts, &response.Meta{Cursor: result.Cursor, HasMore: result.HasMore})
}

func (h *Handler) ListExplore(c *gin.Context) {
	var query FeedQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid query parameters"))
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
	response.OKWithMeta(c, result.Posts, &response.Meta{Cursor: result.Cursor, HasMore: result.HasMore})
}


func (h *Handler) Update(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		response.Error(c, apperrors.Unauthorized("Authentication required"))
		return
	}
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid post ID"))
		return
	}
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperrors.BadRequest("Invalid request: "+err.Error()))
		return
	}
	result, err := h.service.Update(c.Request.Context(), postID, userID, &req)
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
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid post ID"))
		return
	}
	role := c.GetString(middleware.ContextKeyUserRole)
	if err := h.service.Delete(c.Request.Context(), postID, userID, role); err != nil {
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
	postID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperrors.BadRequest("Invalid post ID"))
		return
	}

	post, err := h.service.GetByID(c.Request.Context(), postID)
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
