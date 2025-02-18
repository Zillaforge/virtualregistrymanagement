package tables

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MemberAcl define the database `member_acl` table schema
type MemberAcl struct {
	ID     string `gorm:"not null;primary_key;type:varchar(64)"`
	TagID  string `gorm:"type:varchar(36);not null;uniqueIndex:idx_tid_uid"`
	UserID string `gorm:"type:varchar(36);not null;uniqueIndex:idx_tid_uid"`

	Tag Tag `gorm:"foreignKey:TagID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	_   struct{}
}

// BeforeCreate ...
func (u *MemberAcl) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		for {
			u.ID = uuid.Must(uuid.NewRandom()).String()
			var count int64
			if err = tx.Model(&MemberAcl{}).Where("id = ?", u.ID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				break
			}
		}
	}
	return nil
}
