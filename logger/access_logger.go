package logger

import (
	cnt "VirtualRegistryManagement/constants"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
	"pegasus-cloud.com/aes/toolkits/mviper"
)

type TriggerType string

const (
	API    TriggerType = "api"
	GRPC   TriggerType = "grpc"
	Manual TriggerType = "manual"

	RFC3339Milli = "2006-01-02T15:04:05.999Z07:00"
)

var _accessLogger *lumberjack.Logger

func NewAccessParams() *AccessFormatterParams {
	p := new(AccessFormatterParams)
	p.CloudInfraLogType = "json"
	p.Action.Time = time.Now().Format(RFC3339Milli)
	p.Service.ID = cnt.GetAccessLoggerServiceIDStr()
	p.Service.Name = cnt.Kind
	p.Meta.Location = mviper.GetString("location_id")
	p.Meta.HostName = mviper.GetString("host_id")
	p.Meta.AvailabilityDistrict = mviper.GetString("VirtualRegistryManagement.scopes.availability_district")
	return p
}

func InitAccessLogger(fileName string) {
	_accessLogger = &lumberjack.Logger{
		Filename:   filepath.Join(mviper.GetString("VirtualRegistryManagement.logger.access_log.path"), fileName),
		MaxSize:    mviper.GetInt("VirtualRegistryManagement.logger.access_log.max_size"),    // MB
		MaxBackups: mviper.GetInt("VirtualRegistryManagement.logger.access_log.max_backups"), // how many backup for log file.
		MaxAge:     mviper.GetInt("VirtualRegistryManagement.logger.access_log.max_age"),
		Compress:   mviper.GetBool("VirtualRegistryManagement.logger.access_log.compress"),
	}
}

func writer(p *AccessFormatterParams) {
	accessMarshal, _ := json.Marshal(p)
	go _accessLogger.Write([]byte(fmt.Sprintf("%s\n", accessMarshal)))
}
