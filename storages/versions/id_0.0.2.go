package versions

import (
	"VirtualRegistryManagement/storages/tables"

	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func getID002Migrate() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "0.0.2",
		Migrate: func(tx *gorm.DB) error {
			for _, todo := range []func(tx *gorm.DB) (err error){
				createExportTable,
			} {
				if err := todo(tx); err != nil {
					return err
				}
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			for _, todo := range []func(tx *gorm.DB) (err error){
				dropExportTable,
			} {
				if err := todo(tx); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

// export
func createExportTable(tx *gorm.DB) (err error) {
	if !tx.Migrator().HasTable("export") {
		return tx.Set("gorm:table_options", "ENGINE=InnoDB  CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci").
			Migrator().CreateTable(tables.Export{})
	}
	return nil
}

func dropExportTable(tx *gorm.DB) (err error) {
	if tx.Migrator().HasTable("export") {
		return tx.Migrator().DropTable(tables.Export{})
	}
	return nil
}
