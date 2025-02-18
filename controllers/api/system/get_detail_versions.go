package system

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/utility"
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtil "pegasus-cloud.com/aes/toolkits/utilities"
)

type DetailVersionOutput struct {
	Version        string   `json:"version"`
	GolangVersion  string   `json:"golangVersion"`
	Modules        []string `json:"modules"`
	EventPublishes []string `json:"eventPublishes"`
	AuthChains     []string `json:"authChains"`
}

var (
	prefixes []string = []string{
		"pegasus-cloud.com",
	}
)

func GetDetailVersions(c *gin.Context) {
	f := tracer.StartWithGinContext(
		c,
		tkUtil.NameOfFunction().String(),
	)

	output := &DetailVersionOutput{
		Modules:        make([]string, 0),
		EventPublishes: make([]string, 0),
		AuthChains:     make([]string, 0),
	}

	var statusCode int = http.StatusOK
	defer f(tracer.Attributes{
		"output":     &output,
		"statusCode": &statusCode,
	})

	ScanVersion(output)

	utility.ResponseWithType(c, statusCode, output)
}

func ScanVersion(out *DetailVersionOutput) {
	out.GolangVersion = runtime.Version()
	out.Version = cnt.Version

	modulesVersion(out)
}

func modulesVersion(out *DetailVersionOutput) {
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range bi.Deps {
			for _, prefix := range prefixes {
				if strings.HasPrefix(dep.Path, prefix) {
					out.Modules = append(out.Modules, fmt.Sprintf("%s: %s", dep.Path, dep.Version))
				}
			}
		}
	}
}
