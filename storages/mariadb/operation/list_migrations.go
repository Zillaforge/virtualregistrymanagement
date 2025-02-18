package operation

import "VirtualRegistryManagement/storages/common"

// ListMigrations ...
func (o *Operation) ListMigrations() (migrations []*common.Migration, err error) {
	if err = o.conn.Find(&migrations).Error; err != nil {
		return nil, err
	}
	return migrations, nil
}
