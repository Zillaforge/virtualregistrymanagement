package grpc

import (
	cnt "VirtualRegistryManagement/constants"
	"fmt"
	"os"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/eventpublishpluginclient/epp"
	eppUtil "pegasus-cloud.com/aes/eventpublishpluginclient/utility"
)

type newGRPCConnInput struct {
	conn       int
	socketPath string
	_          struct{}
}

func newGRPCConn(in *newGRPCConnInput) (conn *epp.PoolHandler) {
	newInput := epp.PoolProvider{
		Mode: epp.UnixMode,
		UnixProvider: epp.UnixProvider{
			SocketPath: in.socketPath,
			ConnCount:  in.conn,
		},
		RouteResponseType: eppUtil.JSON,
		Timeout:           2,
	}
	hdr, err := epp.New(newInput)
	if err != nil {
		zap.L().With(
			zap.String(Plugin, "epp.New(....)"),
			zap.Any("input", newInput),
		).Error(err.Error())
		fmt.Println(cnt.PluginInternalServerErrMsg)
		os.Exit(1)
	}
	return hdr
}
