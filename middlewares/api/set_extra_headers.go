package api

import (
	cnt "VirtualRegistryManagement/constants"

	"github.com/gin-gonic/gin"
	"github.com/Zillaforge/toolkits/mviper"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtil "github.com/Zillaforge/toolkits/utilities"
)

// SetExtraHeaders ...SetExtraHeaders
func SetExtraHeaders(c *gin.Context) {
	f := tracer.StartWithGinContext(
		c,
		tkUtil.NameOfFunction().String(),
	)
	defer f(
		tracer.Attributes{
			cnt.HdrHostID:     mviper.GetString("host_id"),
			cnt.HdrLocationID: mviper.GetString("location_id"),
			cnt.HdrVersionID:  cnt.Version,
		},
	)
	c.Set(cnt.HdrLocationID, mviper.GetString("location_id"))
	c.Set(cnt.HdrHostID, mviper.GetString("host_id"))

	c.Header(cnt.HdrLocationID, mviper.GetString("location_id"))
	c.Header(cnt.HdrHostID, mviper.GetString("host_id"))
	c.Header(cnt.HdrVersionID, cnt.Version)
	c.Next()
}
