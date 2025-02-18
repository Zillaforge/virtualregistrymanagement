package utility

func StringInSlice(findData string, list []string) bool {
	for _, b := range list {
		if b == findData {
			return true
		}
	}
	return false
}
