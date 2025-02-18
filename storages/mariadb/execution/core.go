package execution

import (
	"database/sql"

	"VirtualRegistryManagement/storages/mariadb/common"
	"VirtualRegistryManagement/storages/mariadb/operation"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Execution ...
type Execution struct {
	mg *gormigrate.Gormigrate
	op *operation.Operation

	conn *gorm.DB
	db   *sql.DB
	conf *common.ConnectionConfig
}

// SetConfig ...
func (e *Execution) SetConfig(conf *common.ConnectionConfig) {
	e.conf = conf
}

// Set ...
func (e *Execution) Set(mg *gormigrate.Gormigrate, op *operation.Operation, conn *gorm.DB) {
	e.mg = mg
	e.op = op
	e.conn = conn
	e.db, _ = conn.DB()
}
