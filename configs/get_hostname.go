package configs

import (
	"os"
)

func getHostname() string {
	if name, err := os.Hostname(); err == nil {
		return name
	}
	return "Unknown"
}
