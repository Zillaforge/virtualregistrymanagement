package logger

import (
	cnt "VirtualRegistryManagement/constants"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	pbac "github.com/Zillaforge/toolkits/pbac/gin"
)

type AccessFormatterParams struct {
	// CloudInfraLogType is constant value (json). DO NOT MODIFY
	CloudInfraLogType string `json:"cloudinfra_log_type"`
	Service           struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"service"`

	Resource struct {
		// ID is constant value (empty string). DO NOT MODIFY
		ID string `json:"id"`
		// Name is constant value (empty string). DO NOT MODIFY
		Name string `json:"name"`
	} `json:"resource"`

	Project struct {
		ID string `json:"id"`
	} `json:"project"`

	Source struct {
		IP string `json:"ip"`
	} `json:"source"`

	User struct {
		ID string `json:"id"`
	} `json:"user"`

	Action struct {
		// Time is operation time. This is auto generation. DO NOT MODIFY
		Time string `json:"time"`
		// Message is constant value (empty string). DO NOT MODIFY
		Message string `json:"message"`
		ID      string `json:"id"`
		Name    string `json:"name"`
	} `json:"action"`

	Meta MetaFormatterParams `json:"meta"`

	// Simulation Access API Token (SAAT, 模擬分身)
	SAATUser struct {
		ID string `json:"id"`
	} `json:"saatUser"`
}

type MetaFormatterParams struct {
	TriggerType          TriggerType `json:"triggerType"`
	Method               string      `json:"method"`
	StatusCode           int         `json:"statusCode"`
	Path                 string      `json:"path"`
	Latency              string      `json:"latency"`
	RequestID            string      `json:"requestID"`
	HostName             string      `json:"hostName"`
	Location             string      `json:"location"`
	AvailabilityDistrict string      `json:"ad"`
}

func (p *AccessFormatterParams) SetUserID(v string) *AccessFormatterParams {
	p.User.ID = v
	return p
}

func (p *AccessFormatterParams) SetProjectID(v string) *AccessFormatterParams {
	p.Project.ID = v
	return p
}

func (p *AccessFormatterParams) SetProjectIDByContext(c *gin.Context) *AccessFormatterParams {
	if v := c.GetString(cnt.CtxProjectID); v != "" {
		p.Project.ID = v
		return p
	}

	if v := c.Param(cnt.ParamProjectID); v != "" {
		p.Project.ID = v
	}
	return p
}

func (p *AccessFormatterParams) SetSourceIP(v string) *AccessFormatterParams {
	p.Source.IP = v
	return p
}

func (p *AccessFormatterParams) SetAccessLoggerInfo(info cnt.AccessLoggerInfo) *AccessFormatterParams {
	p.Action.ID = strconv.Itoa(info.ID)
	p.Action.Name = info.Name
	return p
}

func (p *AccessFormatterParams) SetAccessLoggerInfoByPBAC(c *gin.Context) *AccessFormatterParams {
	hasRouter, entity := pbac.Actions.Checker(c)
	if !hasRouter {
		zap.L().With(
			zap.String("Logger", "logger.SetAccessLoggerInfoByPBAC(...)"),
			zap.Any("hasRouter", hasRouter)).Warn("cannot find router on PBAC")
		return p
	}

	info := cnt.GetAccessLoggerInfo(entity.Action)

	if info == nil {
		zap.L().With(zap.Any("entity", entity)).Warn("cannot find Entity of PBAC")
		return p
	}

	p.Action.ID = strconv.Itoa(info.ID)
	p.Action.Name = info.Name

	return p
}

/*
	Meta Operations
*/

func (p *AccessFormatterParams) SetTriggerType(v TriggerType) *AccessFormatterParams {
	p.Meta.TriggerType = v
	return p
}

func (p *AccessFormatterParams) SetRequestID(v string) *AccessFormatterParams {
	p.Meta.RequestID = v
	return p
}

// 模擬分身的實際 UserID
func (p *AccessFormatterParams) SetSAATUserIDByContext(c *gin.Context) *AccessFormatterParams {
	p.SAATUser.ID = c.GetString(cnt.CtxSAATUserID)
	return p
}

func (p *AccessFormatterParams) Writer() {
	writer(p)
}
