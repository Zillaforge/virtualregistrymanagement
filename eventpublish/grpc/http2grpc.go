package grpc

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/Zillaforge/eventpublishpluginclient/pb"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtil "github.com/Zillaforge/toolkits/utilities"
)

// HttpRouterToGrpc is unAuthentication api
func (c *core) http2grpc(gc *gin.Context) {
	f := tracer.StartWithGinContext(
		gc,
		tkUtil.NameOfFunction().String(),
	)

	var err error
	var statusCode int = http.StatusOK
	var message string
	var errCode int
	var data interface{}

	defer f(tracer.Attributes{
		"err":        &err,
		"statusCode": &statusCode,
		"message":    &message,
		"errCode":    &errCode,
	})

	md := make(map[string]string)
	for key, values := range gc.Request.Header {
		md[key] = strings.Join(values, ",")
	}

	ginParams := make(map[string]string)
	for _, param := range gc.Params {
		ginParams[param.Key] = param.Value
	}

	ginQuery := make(map[string]string)
	for key, param := range gc.Request.URL.Query() {
		ginQuery[key] = strings.Join(param, ",")
	}

	body, err := ioutil.ReadAll(gc.Request.Body)
	if err != nil {
		statusCode = http.StatusBadRequest
		errCode = cnt.ControllerInternalServerErr.Code()
		message = cnt.ControllerInternalServerErr.Message()
		zap.L().With(
			zap.String(cnt.Controller, "ioutil.ReadAll(...)"),
			zap.String(cnt.RequestID, gc.GetString(cnt.RequestID)),
			zap.Any("obj", gc.Request.Body),
		).Error(err.Error())
		utility.ResponseWithType(gc, statusCode, &utility.ErrResponse{
			ErrorCode: errCode,
			Message:   message,
		})
		return
	}

	prefixRemove := "/iam/api/" + cnt.APIVersion + c.name

	requestInfoInput := &pb.HttpRequestInfo{
		Method:  gc.Request.Method,
		Headers: md,
		Body:    body,
		Path:    strings.TrimPrefix(gc.FullPath(), prefixRemove),
		Params:  ginParams,
		Query:   ginQuery,
	}

	output, err := c.handler.EnableHttpRouter(requestInfoInput, gc)
	if err != nil {
		statusCode = http.StatusInternalServerError
		errCode = cnt.ControllerInternalServerErr.Code()
		message = cnt.ControllerInternalServerErr.Message()
		zap.L().With(
			zap.String(cnt.Controller, "poolHandler.EnableHttpRouter(...)"),
			zap.String(cnt.RequestID, gc.GetString(cnt.RequestID)),
			zap.Any("requestInfoInput", requestInfoInput),
		).Warn(err.Error())
		utility.ResponseWithType(gc, statusCode, &utility.ErrResponse{
			ErrorCode: errCode,
			Message:   message,
		})
		return
	}

	if output.Body != nil {
		if err = json.Unmarshal(output.Body, &data); err != nil {
			statusCode = http.StatusInternalServerError
			errCode = cnt.ControllerInternalServerErr.Code()
			message = cnt.ControllerInternalServerErr.Message()
			zap.L().With(
				zap.String(cnt.Controller, "poolHandler.EnableHttpRouter(...)"),
				zap.String(cnt.RequestID, gc.GetString(cnt.RequestID)),
				zap.Any("requestInfoInput", requestInfoInput),
			).Error(err.Error())
			utility.ResponseWithType(gc, statusCode, &utility.ErrResponse{
				ErrorCode: errCode,
				Message:   message,
			})
			return
		}
	}

	utility.ResponseWithType(gc, int(output.StatusCode), data)
}
