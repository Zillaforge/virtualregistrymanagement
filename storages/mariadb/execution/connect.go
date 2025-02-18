package execution

import (
	"VirtualRegistryManagement/storages/mariadb/common"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// ENDPOINT ...
const (
	endPoint string = "%s:%s@tcp(%s)/%s?charset=utf8&timeout=%ds&parseTime=true&loc=UTC"
)

// Connect ...
func (e *Execution) Connect() (conn *gorm.DB, err error) {
	// check the database exist
	exist, err := e.CheckDatabaseExist(e.conf.DBName)
	if err != nil {
		return
	} else if exist {
		// database already exist
		return connect(e.conf)
	}

	// database does not exist
	// create the specified database
	if err = e.CreateDatabase(e.conf.DBName); err != nil {
		return
	}
	if e.conn, err = connect(e.conf); err != nil {
		return nil, err
	}

	// sync the database
	if err = e.Sync(); err != nil {
		return nil, err
	}
	// init the default data
	if err = e.InitDefaultData(); err != nil {
		return nil, err
	}
	return e.conn, nil
}

func connect(conf *common.ConnectionConfig) (conn *gorm.DB, err error) {
	// generate the endpoint info
	endpoint := fmt.Sprintf(endPoint, conf.Account, conf.Password, conf.Host, conf.DBName, conf.Timeout)
	// connect to database
	conn, err = gorm.Open(mysql.Open(endpoint), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode((logger.Silent)),
	})
	if err != nil {
		return conn, err
	}
	var db *sql.DB
	db, err = conn.DB()
	if err != nil {
		return nil, err
	}
	// setting database connection restriction
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetConnMaxIdleTime(time.Duration(conf.MaxLifetime) * time.Second)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	return conn, nil
}

// Close the database
func close(conn *gorm.DB) {
	if conn != nil {
		if db, _ := conn.DB(); db != nil {
			db.Close()
		}
	}
}
