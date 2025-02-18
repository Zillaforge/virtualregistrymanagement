package system

import (
	cnt "VirtualRegistryManagement/constants"
	util "VirtualRegistryManagement/utility"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"pegasus-cloud.com/aes/toolkits/mviper"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

type ListLogsOutput []Log

type Log struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	From string `json:"from"`
}

var walkPP = []string{
	"virtualregistrymanagement.logger.system_log.path",
	"virtualregistrymanagement.logger.access_log.path",
}

func ListLogs(c *gin.Context) {
	var (
		err        error
		requestID      = util.MustGetContextRequestID(c)
		funcName       = tkUtils.NameOfFunction().Name()
		statusCode int = http.StatusOK
		output         = make(ListLogsOutput, 0)
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	output = append(output, getLogFiles(requestID)...)

	util.ResponseWithType(c, statusCode, output)
}

func getLogFiles(requestID string) (output []Log) {
	output = make([]Log, 0)
	var deDuplicate = map[string]bool{}
	for _, pp := range walkPP {
		innErr := filepath.Walk(mviper.GetString(pp), func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if _, ok := deDuplicate[info.Name()]; ok {
				return nil
			}

			l := Log{
				Name: info.Name(),
				Size: info.Size(),
				From: pp,
			}
			output = append(output, l)
			deDuplicate[l.Name] = true

			return nil
		})
		if innErr != nil {
			zap.L().With(
				zap.String(cnt.Controller, " filepath.Walk(...)"),
				zap.String(cnt.RequestID, requestID),
			).Warn(innErr.Error())
			continue
		}
	}
	return
}
