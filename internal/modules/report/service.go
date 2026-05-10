package report

import (
	"context"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo *Repository
	log  zerolog.Logger
}

func NewService(repo *Repository, log zerolog.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log.With().Str("module", "report").Logger(),
	}
}

func (s *Service) CreateReport(ctx context.Context, reporterID uuid.UUID, req *CreateReportRequest) (*ReportResponse, error) {
	targetID, err := uuid.Parse(req.TargetID)
	if err != nil {
		return nil, apperrors.BadRequest("Invalid target ID")
	}

	already, err := s.repo.HasUserReported(ctx, reporterID, req.TargetType, targetID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if already {
		return nil, apperrors.Conflict("You have already reported this content")
	}

	now := time.Now()
	report := &models.Report{
		ID:          uuid.New(),
		ReporterID:  reporterID,
		TargetType:  req.TargetType,
		TargetID:    targetID,
		Reason:      req.Reason,
		Description: req.Description,
		Status:      models.ReportStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, report); err != nil {
		return nil, apperrors.Internal(err)
	}

	return s.toResponse(report), nil
}

func (s *Service) ListPending(ctx context.Context, limit, offset int) ([]ReportResponse, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	reports, total, err := s.repo.ListPending(ctx, limit, offset)
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	responses := make([]ReportResponse, 0, len(reports))
	for _, r := range reports {
		responses = append(responses, *s.toResponse(&r))
	}
	return responses, total, nil
}

func (s *Service) ListAll(ctx context.Context, limit, offset int) ([]ReportResponse, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	reports, total, err := s.repo.ListAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	responses := make([]ReportResponse, 0, len(reports))
	for _, r := range reports {
		responses = append(responses, *s.toResponse(&r))
	}
	return responses, total, nil
}

func (s *Service) ReviewReport(ctx context.Context, reportID, reviewerID uuid.UUID, req *ReviewReportRequest) error {
	report, err := s.repo.GetByID(ctx, reportID)
	if err != nil {
		return apperrors.Internal(err)
	}
	if report == nil {
		return apperrors.NotFound("Report")
	}
	return s.repo.UpdateStatus(ctx, reportID, req.Status, reviewerID, req.ReviewNote)
}

func (s *Service) toResponse(r *models.Report) *ReportResponse {
	resp := &ReportResponse{
		ID:          r.ID.String(),
		ReporterID:  r.ReporterID.String(),
		TargetType:  r.TargetType,
		TargetID:    r.TargetID.String(),
		Reason:      r.Reason,
		Description: r.Description,
		Status:      string(r.Status),
		ReviewNote:  r.ReviewNote,
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
	}
	if r.ReviewedBy != nil {
		s := r.ReviewedBy.String()
		resp.ReviewedBy = &s
	}
	return resp
}
