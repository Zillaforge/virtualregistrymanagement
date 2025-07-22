package common

import (
	cnt "VirtualRegistryManagement/constants"
	"encoding/json"
	"time"

	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/memkvdb"
	"github.com/Zillaforge/toolkits/mviper"
)

func RetrieveFromCache(key string, isReadFromCache bool, out interface{}, missCallback func() error, add2Cache func() interface{}) error {
	if !memkvdb.IsEnableCache() {
		return missCallback()
	}
	if !isReadFromCache || !memkvdb.Exist(key) {
		if err := missCallback(); err != nil {
			return err
		}
		if add2Cache != nil && isReadFromCache {
			b, _ := json.Marshal(add2Cache())
			// Resetting TTL
			memkvdb.SetTTLBytes(key, b, time.Duration(mviper.GetInt("VirtualRegistryManagement.scopes.memcache_ttl")))
		}
		return nil
	} else {
		b := memkvdb.GetBytes(key)
		if err := json.Unmarshal(b, out); err != nil {
			return tkErr.New(cnt.AuthUnmarshalFromCacheErr).WithInner(err)
		}
		// Resetting TTL
		memkvdb.SetTTLBytes(key, b, time.Duration(mviper.GetInt("VirtualRegistryManagement.scopes.memcache_ttl")))
	}
	return nil
}
