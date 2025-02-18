package common

// ConnectionConfig ...
type ConnectionConfig struct {
	Account      string
	Password     string
	Host         string
	Timeout      int
	MaxOpenConns int
	MaxLifetime  int
	MaxIdleConns int
	DBName       string
}
