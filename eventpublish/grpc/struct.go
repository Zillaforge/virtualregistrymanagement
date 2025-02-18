package grpc

import (
	"os/exec"

	"pegasus-cloud.com/aes/eventpublishpluginclient/epp"
)

type core struct {
	name string
	t    string

	cmd     *exec.Cmd
	handler *epp.PoolHandler
}
