package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	SnapshotTypeDaily   = "daily"
	SnapshotTypeWeekly  = "weekly"
	SnapshotTypeMonthly = "monthly"
	SnapshotTypeManual  = "manual"

	SnapshotStorageLocal = "local"
	SnapshotStorageS3    = "s3"
	SnapshotStorageB2    = "b2"
	SnapshotStorageNone  = "none"

	SnapshotStatusCompleted = "completed"
	SnapshotStatusFailed    = "failed"
	SnapshotStatusSkipped   = "skipped"
)

type WorldSnapshot struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SnapshotType string    `gorm:"not null;column:snapshot_type" json:"snapshotType"`
	Storage      string    `gorm:"not null" json:"storage"`
	Path         string    `gorm:"not null" json:"path"`
	SizeBytes    int64     `gorm:"not null;default:0;column:size_bytes" json:"sizeBytes"`
	Checksum     string    `json:"checksum,omitempty"`
	Status       string    `gorm:"not null;default:completed" json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (WorldSnapshot) TableName() string { return "world_snapshots" }
