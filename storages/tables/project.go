package tables

// Project define the database `project` table schema
type Project struct {
	ID             string `gorm:"type:varchar(36);not null;primary_key"`
	LimitCount     int64  `gorm:"type:bigint(21);not null"`
	LimitSizeBytes int64  `gorm:"type:bigint(21);not null"`

	_ struct{}
}
