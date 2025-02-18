package operation

import (
	"VirtualRegistryManagement/storages/tables"
	"VirtualRegistryManagement/utility/querydecoder"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var _columnMap map[string]string = make(map[string]string)

// Operation ...
type Operation struct {
	conn *gorm.DB
	db   *sql.DB
}

// Set ...
func (o *Operation) Set(conn *gorm.DB, db *sql.DB) {
	o.conn = conn
	o.db = db
	for _, table := range []interface{}{
		// TODO: define tables
		tables.Project{},
		tables.Repository{},
		tables.Tag{},
		tables.MemberAcl{},
		tables.ProjectAcl{},
		tables.Registry{},
		tables.Export{},
	} {
		setColumnMap(table)
	}
}

func setColumnMap(table interface{}) (err error) {
	t := reflect.TypeOf(table)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Name != "_" {
			_columnMap[t.Field(i).Name] = schema.NamingStrategy{}.ColumnName("", t.Field(i).Name)
		}
	}
	return nil
}

func queryConversion(input interface{}) (output map[string]interface{}) {
	t, v := reflect.TypeOf(input), reflect.ValueOf(input)
	output = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			switch v.Field(i).Interface().(type) {
			case querydecoder.Query:
				continue
			}
			output[_columnMap[t.Field(i).Name]] = v.Field(i).Interface()
		}
	}
	return output
}

// GetTableName ...
func GetTableName(table interface{}) (tableName string) {
	names := strings.Split(reflect.TypeOf(table).String(), ".")
	name := names[len(names)-1]
	return schema.NamingStrategy{
		SingularTable: true,
	}.TableName(name)
}

func whereCascade(tx *gorm.DB, input querydecoder.QueryInterface) *gorm.DB {
	querydecoder.Check(input)
	t, v := reflect.TypeOf(input), reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsZero() {
			continue
		}

		prefix := t.Field(i).Tag.Get("prefix")
		switch v.Field(i).Interface().(type) {
		case querydecoder.Query:
			if columnName := t.Field(i).Tag.Get("label"); columnName != "" {
				if prefix != "" {
					columnName = fmt.Sprintf("`%s`.`%s`", prefix, columnName)
				}
				for _, label := range input.GetLabels() {
					tx = tx.Where(fmt.Sprintf("%s LIKE '%%%s%%'", columnName, label))
				}
			}
			continue
		}
		typ := v.Field(i).Kind()
		if typ == reflect.Ptr {
			typ = v.Field(i).Elem().Kind()
		}
		columnName := _columnMap[t.Field(i).Name]
		if prefix != "" {
			columnName = fmt.Sprintf("`%s`.`%s`", prefix, columnName)
		}
		whereKey := strings.ToLower(t.Field(i).Tag.Get("where"))
		switch typ {
		case reflect.Slice:
			where := []string{}
			for _, condition := range input.GetWhere(whereKey) {
				where = append(where, fmt.Sprintf("'%s'", condition.Value))
			}
			if len(where) != 0 {
				tx = tx.Where(fmt.Sprintf("%s IN (%s)", columnName, strings.Join(where, ",")))
			}
		default:
			for _, condition := range input.GetWhere(whereKey) {
				// for query null column
				if condition.Value == "nil" {
					tx = tx.Where(fmt.Sprintf("%s IS Null", columnName))
				} else {
					tx = tx.Where(fmt.Sprintf("%s %s ?",
						columnName,
						symbolConvert(condition.Operator)),
						condition.Value)
				}
			}
		}
	}
	return tx
}

func symbolConvert(symbol string) string {
	switch symbol {
	case "=", "==":
		return "="
	case "!=", "<>":
		return "<>"
	case ">", "<", "<=", ">=":
		return symbol
	default:
		return "="
	}
}
