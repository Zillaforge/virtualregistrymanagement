package execution

import (
	"VirtualRegistryManagement/storages/versions"

	"VirtualRegistryManagement/storages/tables"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Sync ...
func (e *Execution) Sync() (err error) {
	// init the migrate
	e.mg = gormigrate.New(e.conn, &gormigrate.Options{TableName: tables.Migrate}, versions.Get())
	// init the table
	e.mg.InitSchema(func(tx *gorm.DB) (err error) {
		if err := tx.Set("gorm:table_options", "ENGINE=InnoDB  CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci").Migrator().CreateTable(
			// TODO: define tables
			tables.Project{},
			tables.Repository{},
			tables.Tag{},
			tables.MemberAcl{},
			tables.ProjectAcl{},
			tables.Export{},
		); err != nil {
			return err
		}
		return nil
	})
	if err = e.mg.Migrate(); err != nil {
		return err
	}

	// create view
	for _, d := range []string{
		// TODO: define views
		tables.RegistryView(),
	} {
		if err := e.conn.Exec(d).Error; err != nil {
			return err
		}
	}
	return nil
}
