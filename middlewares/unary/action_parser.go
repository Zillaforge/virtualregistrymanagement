package unary

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/eventpublish"
	"VirtualRegistryManagement/utility"
	"context"
	"fmt"
	"path"
	"regexp"

	"google.golang.org/grpc"
	"pegasus-cloud.com/aes/toolkits/tracer"
)

// ActionParser ...
func ActionParser() grpc.UnaryServerInterceptor {
	// closure function
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md := map[string]string{}
		md[tracer.RequestID] = utility.MustGetContextRequestID(ctx)
		wvCtx := context.WithValue(ctx, cnt.CtxWithExtraVal{}, md)

		//resp, err := handler(ctx, req)
		resp, err := handler(wvCtx, req)
		if err != nil {
			return resp, err
		}

		// 如果 info.FullMethod 為 filter 中的值，則忽略
		if !utility.StringInSlice(parseRPCMethod(info.FullMethod), []string{}) {
			eventpublish.GetBus().Publish(cnt.ReconcileKey, genGRPCAction(path.Base(info.FullMethod)), md, req, resp)
		}
		return resp, err
	}
}

// parseRPCMethod will return service name, i.e., /package.service/method
func parseRPCMethod(fullMethod string) string {
	re := regexp.MustCompile(`/(?P<package>[a-zA-Z0-9-_]+).(?P<service>[a-zA-Z0-9-_]+)/(?P<method>[a-zA-Z0-9-_]+)`)
	matches := re.FindStringSubmatch(fullMethod)
	return matches[re.SubexpIndex("service")]
}

func genGRPCAction(method string) string {
	return fmt.Sprintf("grpc:%s", method)
}
