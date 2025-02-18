package execution

import (
	"VirtualRegistryManagement/storages/common"
	v "VirtualRegistryManagement/storages/versions"
)

const defaultSchema string = "SCHEMA_INIT"

// ListMigrationMap returns all of migration versions which can be used from registered list
// and the value of current is which version is using in currently.
func (e *Execution) ListMigrationMap() (current string, versions []string, err error) {
	current, versions = defaultSchema, []string{defaultSchema}
	// get database migration table data
	var records []*common.Migration
	records, err = e.op.ListMigrations()
	if err != nil {
		return "", nil, err
	}
	// get migration map data
	for _, migrateMap := range v.Get() {
		versions = append(versions, migrateMap.ID)
	}

	// if database migration only one record,it's the SCHEMA_INIT
	if len(records) == 1 {
		return defaultSchema, versions, nil
	}

	// find the database current migration id by migration map
	for i := len(versions) - 1; i >= 0; i-- {
		for _, r := range records {
			if r.ID == versions[i] {
				if r.ID == defaultSchema {
					return "", versions, nil
				}
				return r.ID, versions, nil
			}
		}
	}
	return "", versions, nil
}
