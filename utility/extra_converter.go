package utility

import (
	"encoding/json"
)

func Extra2Bytes(v map[string]interface{}) []byte {
	bytes, err := json.Marshal(v)
	if err != nil {
		return []byte{0x7B, 0x7D}
	}
	return bytes
}

func Bytes2Extra(b []byte) map[string]interface{} {
	v := new(map[string]interface{})
	json.Unmarshal(b, v)
	return *v
}
