package execution

import (
	"VirtualRegistryManagement/storages/mariadb/common"
	"errors"
	"fmt"

	"pegasus-cloud.com/aes/toolkits/mviper"
)

const (
	createDBCmd string = "CREATE DATABASE %s"
)

// CreateDatabase create the specified database
func (e *Execution) CreateDatabase(name string) (err error) {
	// init the createDatabaseConf
	conf := &common.ConnectionConfig{
		Account:  mviper.GetString("storage.mariadb.account"),
		Password: mviper.GetString("storage.mariadb.password"),
		Host:     mviper.GetString("storage.mariadb.host"),
		DBName:   "",
		Timeout:  mviper.GetInt("storage.mariadb.timeout"),
	}
	// connect to database
	conn, err := connect(conf)
	if err != nil {
		return
	}
	defer close(conn)
	// create the specified database
	if err = conn.Exec(fmt.Sprintf(createDBCmd, name)).Error; err != nil {
		return
	}
	// check database already created
	exist, err := e.CheckDatabaseExist(name)
	if err != nil {
		return
	}
	// database not created
	if !exist {
		return errors.New("create database fail")
	}
	return
}
