package execution

import "VirtualRegistryManagement/storages/common"

// ListMigrationHistory ...
func (e *Execution) ListMigrationHistory() (versions []*string, err error) {
	var records []*common.Migration
	records, err = e.op.ListMigrations()
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		versions = append(versions, &record.ID)
	}
	return versions, nil
}
