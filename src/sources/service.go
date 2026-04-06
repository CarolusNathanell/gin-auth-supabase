package sources

import (
	"context"
	"errors"
	"gin-auth-supabase/src/db"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	q *db.Queries
}

func NewService(q *db.Queries) *Service {
	return &Service{q: q}
}

func (s *Service) Add(ctx context.Context, req SourcesAdd) (*db.Source, error) {
	switch req.Type {
	case "MP4":
		if !strings.HasSuffix(req.Url, ".mp4") {
			return nil, errors.New("MP4 type requires a link ending in .mp4")
		}
	case "RTSP":
		if !strings.HasPrefix(req.Url, "rtsp://") {
			return nil, errors.New("RTSP type requires a link starting with rtsp://")
		}
	case "Webcam":
		if !strings.HasPrefix(req.Url, "http://") && !strings.HasPrefix(req.Url, "https://") {
			return nil, errors.New("Webcam streams must be HTTP or HTTPS")
		}
	default:
		return nil, errors.New("Invalid source type")
	}

	source, err := s.q.CreateSource(ctx, db.CreateSourceParams{
		Name:       req.Name,
		Type:       db.Sourcetype(req.Type),
		Url:        req.Url,
		FpsTarget:  req.FpsTarget,
		Resolution: req.Resolution,
		Status:     *req.Status,
	})
	return &source, err
}

func (s *Service) Request(ctx context.Context) (*[]db.Source, error) {
	sources, err := s.q.GetSources(ctx)
	if err != nil {
		return nil, errors.New("Source not found")
	}

	return &sources, nil
}

func (s *Service) RequestByID(ctx context.Context, sourceId uuid.UUID) (*db.Source, error) {
	source, err := s.q.GetSourceByID(ctx, sourceId)
	if err != nil {
		return nil, errors.New("Source not found")
	}

	return &source, nil
}

func (s *Service) UpdateById(ctx context.Context, req SourcesUpdate, sourceId uuid.UUID) (*db.Source, error) {
	switch req.Type {
	case "MP4":
		if !strings.HasSuffix(req.Url, ".mp4") {
			return nil, errors.New("MP4 type requires a link ending in .mp4")
		}
	case "RTSP":
		if !strings.HasPrefix(req.Url, "rtsp://") {
			return nil, errors.New("RTSP type requires a link starting with rtsp://")
		}
	case "Webcam":
		if !strings.HasPrefix(req.Url, "http://") && !strings.HasPrefix(req.Url, "https://") {
			return nil, errors.New("Webcam streams must be HTTP or HTTPS")
		}
	default:
		return nil, errors.New("Invalid source type")
	}

	source, err := s.q.UpdateSource(ctx, db.UpdateSourceParams{
		ID:         sourceId,
		Name:       req.Name,
		Type:       db.Sourcetype(req.Type),
		Url:        req.Url,
		FpsTarget:  req.FpsTarget,
		Resolution: req.Resolution,
		Status:     *req.Status,
	})
	return &source, err
}

func (s *Service) DeleteById(ctx context.Context, sourceId uuid.UUID) (*db.Source, error) {
	source, err := s.q.DeleteSource(ctx, sourceId)
	if err != nil {
		return nil, errors.New("Source not found")
	}

	return &source, nil
}
