package execution

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
)

// AutoMigration ...
func (e *Execution) AutoMigration() (err error) {
	zap.L().Info(fmt.Sprintf("start auto migration database"))
	// get the map
	current, versionMap, err := e.ListMigrationMap()
	if len(versionMap) == 1 {
		return nil
	}
	if current == "" {
		return errors.New("can't find database current version")
	}
	if err := e.Migrate(versionMap[len(versionMap)-1]); err != nil {
		return err
	}
	return nil
}
