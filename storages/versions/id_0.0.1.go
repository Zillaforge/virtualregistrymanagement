package versions

import (
	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type id001Migration struct {
	Protect bool `gorm:"not null;type:TINYINT(1);default:0"`
}

func getID001Migrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "0.0.1",
		Migrate: func(tx *gorm.DB) error {
			for _, todo := range []func(tx *gorm.DB) (err error){
				addColumn001RepositoryProtect,
				addColumn001TagProtect,
			} {
				if err := todo(tx); err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			for _, todo := range []func(tx *gorm.DB) (err error){
				dropColumn001RepositoryProtect,
				dropColumn001TagProtect,
			} {
				if err := todo(tx); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func addColumn001RepositoryProtect(tx *gorm.DB) (err error) {
	return tx.Table("repository").Migrator().AddColumn(&id001Migration{}, "protect")
}

func dropColumn001RepositoryProtect(tx *gorm.DB) (err error) {
	return tx.Table("repository").Migrator().DropColumn(&id001Migration{}, "protect")
}

func addColumn001TagProtect(tx *gorm.DB) (err error) {
	return tx.Table("tag").Migrator().AddColumn(&id001Migration{}, "protect")
}

func dropColumn001TagProtect(tx *gorm.DB) (err error) {
	return tx.Table("tag").Migrator().DropColumn(&id001Migration{}, "protect")
}
