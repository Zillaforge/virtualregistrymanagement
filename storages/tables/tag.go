package tables

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Tag define the database `tag` table schema
type Tag struct {
	ID              string         `gorm:"type:varchar(36);not null;primary_key"`
	Name            string         `gorm:"type:varchar(64);not null;uniqueIndex:idx_name_rid"`
	RepositoryID    string         `gorm:"type:varchar(36);not null;uniqueIndex:idx_name_rid"`
	ReferenceTarget string         `gorm:"type:varchar(36)"`
	Type            string         `gorm:"type:varchar(16);not null"`
	Size            uint64         `gorm:"type:bigint(21) unsigned;default:0"`
	Status          string         `gorm:"type:varchar(16);not null"`
	Extra           datatypes.JSON `gorm:"check:json_valid(Extra);default:'{}'"`
	Protect         bool           `gorm:"not null;type:TINYINT(1);default:0"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Repository Repository `gorm:"foreignKey:RepositoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	_          struct{}
}

// BeforeCreate ...
func (d *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		for {
			d.ID = uuid.Must(uuid.NewRandom()).String()
			var count int64
			if err = tx.Model(&Tag{}).Where("id = ?", d.ID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				break
			}
		}
	}
	return nil
}
