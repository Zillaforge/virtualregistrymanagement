package execution

const (
	latest string = "latest"
	oldest string = "oldest"
)

// Migrate ...
func (e *Execution) Migrate(version string) (err error) {
	switch version {
	case latest:
		return e.mg.Migrate()
	default:
		return e.mg.MigrateTo(version)
	}
}
