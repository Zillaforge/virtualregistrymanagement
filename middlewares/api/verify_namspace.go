package api

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack"
	"VirtualRegistryManagement/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

// VerifyNamespaceInput ...
type VerifyNamespaceInput struct {
	Namespace string
	_         struct{}
}

// VerifyOpstkNamespace ...
func VerifyOpstkNamespace(c *gin.Context) {
	var (
		funcName = tkUtils.NameOfFunction().Name()
		// requestID  = utility.MustGetContextRequestID(c)
		statusCode = http.StatusOK
		err        error
		input      = &VerifyNamespaceInput{}
	)

	f := tracer.StartWithGinContext(c, funcName)

	defer f(tracer.Attributes{
		"input":      input,
		"err":        &err,
		"statusCode": &statusCode,
	})

	input.Namespace = c.GetHeader(cnt.HdrNamespace)

	if !openstack.NamespaceIsLegal(input.Namespace) {
		statusCode = http.StatusForbidden
		err = tkErr.New(cnt.MidPermissionDeniedErr)
		utility.ResponseWithType(c, statusCode, err)
		return
	}

	c.Set(cnt.CtxNamespace, input.Namespace)
	c.Next()
}
