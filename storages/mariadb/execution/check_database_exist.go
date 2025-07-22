package execution

import (
	"VirtualRegistryManagement/storages/mariadb/common"

	"github.com/Zillaforge/toolkits/mviper"
)

const (
	schemaDBName string = "information_schema"
	schemaTable  string = "schemata"
)

// CheckDatabaseExist check the database exist
func (e *Execution) CheckDatabaseExist(name string) (exist bool, err error) {
	// init the checkDatabaseConf
	conf := &common.ConnectionConfig{
		Account:  mviper.GetString("storage.mariadb.account"),
		Password: mviper.GetString("storage.mariadb.password"),
		Host:     mviper.GetString("storage.mariadb.host"),
		DBName:   schemaDBName,
		Timeout:  mviper.GetInt("storage.mariadb.timeout"),
	}
	// connect to database
	conn, err := connect(conf)
	if err != nil {
		return
	}
	defer close(conn)
	// count the specified database name
	var count int64
	if err = conn.Table(schemaTable).Where("schema_name = ?", name).Count(&count).Error; err != nil {
		return
	}
	// database name is exist in schemata
	if count > 0 {
		return true, nil
	}
	return
}
