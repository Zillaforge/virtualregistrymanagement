package system

import (
	"VirtualRegistryManagement/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPlainTextVersion(c *gin.Context) {
	c.String(http.StatusOK, constants.Version)
}
