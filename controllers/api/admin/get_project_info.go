package admin

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

type GetProjectInfoInput struct {
	ProjectID string `json:"-"`
	_         struct{}
}

type GetProjectInfoOutput struct {
	UsedSize       int64 `json:"usedSize"`
	UsedCount      int64 `json:"usedCount"`
	SoftLimitSize  int64 `json:"softLimitSize"`
	SoftLimitCount int64 `json:"softLimitCount"`
	_              struct{}
}

func GetProjectInfo(c *gin.Context) {
	var (
		input  = &GetProjectInfoInput{ProjectID: c.GetString(cnt.CtxProjectID)}
		output = &GetProjectInfoOutput{
			UsedSize:       c.GetInt64(cnt.CtxProjectUsedSize),
			UsedCount:      c.GetInt64(cnt.CtxProjectUsedCount),
			SoftLimitSize:  c.GetInt64(cnt.CtxProjectSoftLimitSize),
			SoftLimitCount: c.GetInt64(cnt.CtxProjectSoftLimitCount),
		}
		err        error
		funcName       = tkUtils.NameOfFunction().Name()
		statusCode int = http.StatusOK
	)

	f := tracer.StartWithGinContext(c, funcName)
	defer f(tracer.Attributes{
		"input":      &input,
		"output":     &output,
		"error":      &err,
		"statusCode": &statusCode,
	})

	utility.ResponseWithType(c, statusCode, output)
}
