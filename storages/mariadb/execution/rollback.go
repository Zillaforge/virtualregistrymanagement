package execution

// Rollback ...
func (e *Execution) Rollback() (err error) {
	return e.mg.RollbackLast()
}
