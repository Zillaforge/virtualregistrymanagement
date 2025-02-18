package common

import "gorm.io/gorm"

// Executions ...
type Executions interface {
	Close() (err error)
	Connect() (conn *gorm.DB, err error)
	IsConnecting() (err error)
	Sync() (err error)
	InitDefaultData() (err error)
	Migrate(version string) (err error)
	Rollback() (err error)
	ListMigrationMap() (current string, versions []string, err error)
	ListMigrationHistory() (versions []*string, err error)
	AutoMigration() (err error)
	CheckDatabaseExist(name string) (exist bool, err error)
	CreateDatabase(name string) (err error)
}
