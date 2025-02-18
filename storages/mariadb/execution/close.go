package execution

// Close the database
func (e *Execution) Close() (err error) {
	if e.conn != nil && e.db != nil {
		return e.db.Close()
	}
	return nil
}
