package services

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type RedisSentinelInput struct {
	Name                                        string
	Hosts                                       []string
	MasterGroupName, Password, SentinelPassword string
}

func UnmarshalRedisSentinel(svcCfg *viper.Viper) {
	InitRedisSentinel(&RedisSentinelInput{
		Name:             svcCfg.GetString("name"),
		Hosts:            svcCfg.GetStringSlice("hosts"),
		MasterGroupName:  svcCfg.GetString("master_group_name"),
		Password:         svcCfg.GetString("password"),
		SentinelPassword: svcCfg.GetString("sentinel_password"),
	})
}

func InitRedisSentinel(input *RedisSentinelInput) {
	conn := newRedisSentinel(input)
	ServiceMap[input.Name] = &Service{
		Kind: _redisSentinelKind,
		Conn: conn,
	}
}

func NewRedisSentinel(input *RedisSentinelInput) (conn *redis.Client) {
	return newRedisSentinel(input)
}

func newRedisSentinel(input *RedisSentinelInput) (conn *redis.Client) {
	return redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       input.MasterGroupName,
		SentinelAddrs:    input.Hosts,
		Password:         input.Password,
		SentinelPassword: input.SentinelPassword,
	})
}
