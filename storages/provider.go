package storages

import (
	"fmt"

	"VirtualRegistryManagement/storages/common"
	"VirtualRegistryManagement/storages/mariadb"
	mariadbCom "VirtualRegistryManagement/storages/mariadb/common"

	"github.com/Zillaforge/toolkits/mviper"
)

const (
	//mariaDBProviderKey define database type for mariadb
	mariaDBProviderKey = "mariadb"
)

// Provider ...
type Provider struct {
	Op   common.Operations
	Exec common.Executions
}

var (
	provider Provider
)

// New new the database instance and return the connect
func New(p string) {
	switch p {
	case mariaDBProviderKey:
		mariaDBProvider, err := mariadb.New(&mariadbCom.ConnectionConfig{
			Account:      mviper.GetString("storage.mariadb.account"),
			Password:     mviper.GetString("storage.mariadb.password"),
			Host:         mviper.GetString("storage.mariadb.host"),
			DBName:       mviper.GetString("storage.mariadb.name"),
			Timeout:      mviper.GetInt("storage.mariadb.timeout"),
			MaxOpenConns: mviper.GetInt("storage.mariadb.max_open_conns"),
			MaxLifetime:  mviper.GetInt("storage.mariadb.conn_max_lifetime"),
			MaxIdleConns: mviper.GetInt("storage.mariadb.max_idle_conns"),
		})
		if err != nil {
			panic(err)
		}
		provider.Exec = &mariaDBProvider.Exec
		provider.Op = &mariaDBProvider.Op
	default:
		panic(fmt.Errorf("not support the %s storage", p))
	}
}

// Replace replaces global provider by p
func Replace(p Provider) {
	provider = p
}

// Use ...
func Use() (op common.Operations) {
	return provider.Op
}

// Exec returns executions interface which is a set of storage kernel functions
func Exec() (exec common.Executions) {
	return provider.Exec
}
