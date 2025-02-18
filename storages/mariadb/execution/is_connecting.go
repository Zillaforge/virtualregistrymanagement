package execution

// IsConnecting ...
func (e *Execution) IsConnecting() (err error) {
	return e.db.Ping()
}
