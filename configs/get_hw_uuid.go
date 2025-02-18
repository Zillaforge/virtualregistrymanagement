package configs

import (
	"VirtualRegistryManagement/constants"
	"os"
)

func getLocation() string {
	if uuidBytes, err := os.ReadFile(constants.ProductUUIDFilePath); err == nil {
		return string(uuidBytes)
	}
	return "Unknown"
}
