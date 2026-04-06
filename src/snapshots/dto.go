package snapshots

import "github.com/google/uuid"

type SnapshotAdd struct {
	ID              int32     `json:"id" binding:"required"`
	SourceID        uuid.UUID `json:"source_id" binding:"required"`
	ImagePath       string    `json:"image_path"`
	HeadCountAtTime int32     `json:"head_count_at_time" binding:"required"`
}
