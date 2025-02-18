package common

import (
	cnt "VirtualRegistryManagement/constants"
	"encoding/json"

	"github.com/gophercloud/gophercloud"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
)

func IsKnownError(err error) (error, bool) {
	if e, ok := err.(gophercloud.ErrUnexpectedResponseCode); !ok {
		return err, false
	} else {
		b := map[string]interface{}{}
		switch e.GetStatusCode() {
		case 413: // Quota: over limit
			if _err := json.Unmarshal(e.Body, &b); _err != nil {
				return err, false
			}
			if _, exist := b["overLimit"]; exist {
				return tkErr.New(cnt.OpenstackExceedAllowedQuotaErr).WithInner(e), true
			}
		}
		return err, false
	}
}
