package grpc

import (
	"os/exec"

	"github.com/Zillaforge/eventpublishpluginclient/epp"
)

type core struct {
	name string
	t    string

	cmd     *exec.Cmd
	handler *epp.PoolHandler
}
