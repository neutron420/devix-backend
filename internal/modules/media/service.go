package media

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"devix-backend/internal/config"
	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/models"
	"devix-backend/internal/pkg/sanitize"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo    *Repository
	storage StorageProvider
	cfg     config.MediaConfig
	log     zerolog.Logger
}

func NewService(repo *Repository, storage StorageProvider, cfg config.MediaConfig, log zerolog.Logger) *Service {
	return &Service{
		repo:    repo,
		storage: storage,
		cfg:     cfg,
		log:     log.With().Str("module", "media").Logger(),
	}
}

func (s *Service) UploadAvatar(ctx context.Context, file multipart.File, header *multipart.FileHeader) (string, error) {
	mimeType, err := s.detectMIME(file)
	if err != nil {
		return "", apperrors.BadRequest("Unable to detect file type")
	}

	if !s.isAllowedImageType(mimeType) {
		return "", apperrors.BadRequest("File type not allowed")
	}

	safeFilename := sanitize.Filename(header.Filename)
	path, err := s.storage.Upload(ctx, file, "avatars", safeFilename)
	if err != nil {
		return "", apperrors.Internal(err)
	}

	return s.storage.GetURL(path), nil
}

func (s *Service) UploadPostMedia(ctx context.Context, postID uuid.UUID, files []*FileUpload) ([]*models.Media, error) {
	var results []*models.Media

	for i, f := range files {
		mimeType, err := s.detectMIME(f.File)
		if err != nil {
			return nil, apperrors.BadRequest(fmt.Sprintf("Unable to detect type of file %d", i+1))
		}

		var fileType models.MediaType
		var maxSize int64

		if s.isAllowedImageType(mimeType) {
			fileType = models.MediaTypeImage
			maxSize = s.cfg.MaxImageSize
		} else if s.isAllowedVideoType(mimeType) {
			fileType = models.MediaTypeVideo
			maxSize = s.cfg.MaxVideoSize
		} else if s.isAllowedDocType(mimeType) {
			fileType = models.MediaTypeDocument
			maxSize = s.cfg.MaxDocSize
		} else {
			return nil, apperrors.BadRequest("File type not allowed")
		}

		// Check size
		if f.Header.Size > maxSize {
			return nil, apperrors.PayloadTooLarge(fmt.Sprintf("File %s is too large", f.Header.Filename))
		}

		// Check limits per post
		count, _ := s.repo.CountByPostAndType(ctx, postID, fileType)
		if fileType == models.MediaTypeImage && int(count) >= s.cfg.MaxImagesPerPost {
			return nil, apperrors.BadRequest("Max images reached")
		}
		if fileType == models.MediaTypeVideo && count >= 1 {
			return nil, apperrors.BadRequest("Only one video allowed per post")
		}
		if fileType == models.MediaTypeDocument && count >= 1 {
			return nil, apperrors.BadRequest("Only one document allowed per post")
		}

		safeFilename := sanitize.Filename(f.Header.Filename)
		path, err := s.storage.Upload(ctx, f.File, fmt.Sprintf("posts/%s", postID), safeFilename)
		if err != nil {
			return nil, apperrors.Internal(err)
		}

		sortOrder, _ := s.repo.GetNextSortOrder(ctx, postID)

		media := &models.Media{
			ID:           uuid.New(),
			PostID:       postID,
			FileURL:      s.storage.GetURL(path),
			FileType:     fileType,
			FileSize:     f.Header.Size,
			MimeType:     mimeType,
			OriginalName: &f.Header.Filename,
			SortOrder:    sortOrder + i,
			CreatedAt:    time.Now(),
		}

		if err := s.repo.Create(ctx, media); err != nil {
			return nil, apperrors.Internal(err)
		}

		results = append(results, media)
	}

	return results, nil
}

func (s *Service) GetPostMedia(ctx context.Context, postID uuid.UUID) ([]models.Media, error) {
	return s.repo.GetByPostID(ctx, postID)
}

func (s *Service) DeletePostMedia(ctx context.Context, postID uuid.UUID) error {
	urls, err := s.repo.DeleteByPostID(ctx, postID)
	if err != nil {
		return err
	}
	for _, u := range urls {
		_ = s.storage.Delete(ctx, u)
	}
	return nil
}

func (s *Service) detectMIME(file multipart.File) (string, error) {
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		return "", err
	}
	_, _ = file.Seek(0, 0)
	return http.DetectContentType(buf[:n]), nil
}

func (s *Service) isAllowedImageType(mimeType string) bool {
	for _, t := range s.cfg.AllowedImageTypes {
		if mimeType == t {
			return true
		}
	}
	return false
}

func (s *Service) isAllowedVideoType(mimeType string) bool {
	for _, t := range s.cfg.AllowedVideoTypes {
		if mimeType == t {
			return true
		}
	}
	return false
}

func (s *Service) isAllowedDocType(mimeType string) bool {
	for _, t := range s.cfg.AllowedDocTypes {
		if mimeType == t {
			return true
		}
	}
	return false
}
