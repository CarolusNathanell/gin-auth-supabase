package sources

import "github.com/google/uuid"

type SourceAdd struct {
	Name       string    `json:"name" binding:"required"`
	Type       string    `json:"type" binding:"required,oneof=MP4 RTSP Webcam Youtube Other"`
	Url        string    `json:"url" binding:"required,url"`
	FpsTarget  int32     `json:"fps_target" binding:"required"`
	Resolution string    `json:"resolution" binding:"required"`
	Status     *bool     `json:"status" binding:"required"`
	UserID     uuid.UUID `json:"user_id" binding:"required"`
}

type SourceUpdate struct {
	Name       string    `json:"name" binding:"required"`
	Type       string    `json:"type" binding:"required,oneof=MP4 RTSP Webcam Youtube Other"`
	Url        string    `json:"url" binding:"required,url"`
	FpsTarget  int32     `json:"fps_target" binding:"required"`
	Resolution string    `json:"resolution" binding:"required"`
	Status     *bool     `json:"status" binding:"required"`
	UserID     uuid.UUID `json:"user_id" binding:"required"`
}

type SourcesResponse struct {
	Sources any `json:"sources"`
}

type SourceDelete struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}
