package tables

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Export define the database `export` table schema
type Export struct {
	ID             string  `gorm:"type:varchar(36);not null;primary_key"`
	RepositoryID   string  `gorm:"type:varchar(36);not null"`
	RepositoryName string  `gorm:"type:varchar(64);not null"`
	TagID          string  `gorm:"type:varchar(36);not null"`
	TagName        string  `gorm:"type:varchar(64);not null"`
	Type           string  `gorm:"type:varchar(16);not null"`
	SnapshotID     *string `gorm:"type:varchar(36)"`
	SnapshotStatus *string `gorm:"type:varchar(16)"`
	VolumeID       *string `gorm:"type:varchar(36)"`
	VolumeStatus   *string `gorm:"type:varchar(16)"`
	ImageID        string  `gorm:"type:varchar(36);not null"`
	ImageStatus    string  `gorm:"type:varchar(16);not null"`
	Filepath       string  `gorm:"type:text;not null"`
	Status         string  `gorm:"type:varchar(16);default:preparing"`
	Creator        string  `gorm:"type:varchar(36);not null"`
	ProjectID      string  `gorm:"type:varchar(36);not null"`
	Namespace      string  `gorm:"type:varchar(16);not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	_              struct{}
}

// BeforeCreate ...
func (d *Export) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		for {
			d.ID = uuid.Must(uuid.NewRandom()).String()
			var count int64
			if err = tx.Model(&Export{}).Where("id = ?", d.ID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				break
			}
		}
	}
	return nil
}
