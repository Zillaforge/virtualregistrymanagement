package api

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
	"pegasus-cloud.com/aes/virtualplatformserviceclient/pb"
	"pegasus-cloud.com/aes/virtualplatformserviceclient/vps"
)

const (
	vrmPrefix = "vrm_extras_"
	osKey = "operating_system"
)

// VerifyOpstkServer ...
func VerifyOpstkServer(c *gin.Context) {
	var (
		funcName = tkUtils.NameOfFunction().Name()
		// requestID  = utility.MustGetContextRequestID(c)
		statusCode = http.StatusOK
		err        error
		input      = &ResourceIDInput{ID: c.Param(cnt.ParamServerID)}
	)

	f := tracer.StartWithGinContext(c, funcName)

	defer f(tracer.Attributes{
		"input":      input,
		"err":        &err,
		"statusCode": &statusCode,
	})

	getOutput, err := vps.Server().Get(&pb.IDInput{ID: input.ID}, c)
	if err != nil {
		statusCode = http.StatusNotFound
		err = tkErr.New(cnt.MidServerNotFoundErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	if c.GetString(cnt.CtxNamespace) != getOutput.Namespace ||
		c.GetString(cnt.CtxProjectID) != getOutput.ProjectID {
		statusCode = http.StatusForbidden
		err = tkErr.New(cnt.MidPermissionDeniedErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	c.Set(cnt.CtxVolumeID, getOutput.RootDiskID)
	c.Set(cnt.CtxCreator, getOutput.UserID)

	// Storing prefix key has "vrm_extras_" and "operating_system" to next level
	vpsMetadata := dictHasPrefixKey(getOutput.Metadatas)
	c.Set(cnt.CtxVPSMetadata, vpsMetadata)
	c.Set(cnt.CtxServerOS, getOutput.Metadatas[osKey])
	c.Next()
}

func dictHasPrefixKey(dict map[string]string) (map[string]string) {
	targetMap := make(map[string]string)
	for key := range dict {
		if strings.HasPrefix(key, vrmPrefix) {
			targetKey := strings.TrimPrefix(key, vrmPrefix)
			targetMap[targetKey] = dict[key]
		}
	}

	return targetMap
}