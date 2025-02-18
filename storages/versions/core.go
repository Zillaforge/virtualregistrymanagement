package versions

import "github.com/go-gormigrate/gormigrate/v2"

// Get the all migration list
func Get() (versions []*gormigrate.Migration) {
	return []*gormigrate.Migration{
		getID001Migrate(),
		getID002Migrate(),
	}
}
