package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/svr"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

const LogPathPrefix = "/api/admin/log"

type LogRoute struct {
	Cfg        *config.Config
	LogService *svr.LogService
}

var LogRouteSet = wire.NewSet(wire.Struct(new(LogRoute), "*"))

// Get
// @Summary get logs
// @Schemes
// @Description get logs
// @Tags Log
// @Accept json
// @Produce json
// @Param        page               query      int        false       "page"               -1
// @Param        pageSize           query      int        false       "pageSize"            5
// @Param        eventType          query      string     false       "eventType"
// @Success 200 {object} apm.QueryResult
// @Router /api/admin/log [Get]
func (l *LogRoute) Get(c *gin.Context) *R {
	option := &credential.LogSearchOption{}

	if err := c.BindQuery(option); err != nil {
		return InternalError(err)
	}
	result, err := l.LogService.GetByOption(option)
	if err != nil {
		return InternalError(err)
	}
	return Ok(result)
}

func (l *LogRoute) Register(r *gin.Engine) {
	r.GET(LogPathPrefix, Handle(l.Get))
}
