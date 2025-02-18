package tables

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository define the database `repository` table schema
type Repository struct {
	ID              string `gorm:"type:varchar(36);not null;primary_key"`
	Name            string `gorm:"type:varchar(64);not null;uniqueIndex:idx_n_c_pid"`
	Namespace       string `gorm:"type:varchar(16);not null;uniqueIndex:idx_n_c_pid"`
	OperatingSystem string `gorm:"type:varchar(16);not null"`
	Description     string `gorm:"type:text;default:NULL"`
	Creator         string `gorm:"type:varchar(36);not null;uniqueIndex:idx_n_c_pid"`
	ProjectID       string `gorm:"type:varchar(36);not null;uniqueIndex:idx_n_c_pid"`
	Protect         bool   `gorm:"not null;type:TINYINT(1);default:0"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Project Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tag     []Tag   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	_       struct{}
}

// BeforeCreate ...
func (d *Repository) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		for {
			d.ID = uuid.Must(uuid.NewRandom()).String()
			var count int64
			if err = tx.Model(&Repository{}).Where("id = ?", d.ID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				break
			}
		}
	}
	return nil
}
