package tables

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProjectAcl define the database `project_acl` table schema
type ProjectAcl struct {
	ID        string  `gorm:"not null;primary_key;type:varchar(64)"`
	TagID     string  `gorm:"type:varchar(36);not null;uniqueIndex:idx_tid_pid"`
	ProjectID *string `gorm:"type:varchar(36);uniqueIndex:idx_tid_pid"`

	Tag     Tag     `gorm:"foreignKey:TagID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Project Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	_       struct{}
}

// BeforeCreate ...
func (u *ProjectAcl) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		for {
			u.ID = uuid.Must(uuid.NewRandom()).String()
			var count int64
			if err = tx.Model(&ProjectAcl{}).Where("id = ?", u.ID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				break
			}
		}
	}
	return nil
}
