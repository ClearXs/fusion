package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/event"
	"cc.allio/fusion/internal/svr"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

const isrPathPrefix = "/api/admin/isr"

type IsrRoute struct {
	Cfg            *config.Config
	SettingService *svr.SettingService
	Isr            *event.IsrEventBus
}

var IsrRouterSet = wire.NewSet(wire.Struct(new(IsrRoute), "*"))

// ActiveAll
// @Summary active isr
// @Schemes
// @Description active isr
// @Tags Isr
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/isr [Post]
func (i *IsrRoute) ActiveAll(c *gin.Context) *R {
	i.Isr.ActiveAll("trigger incremental rendering by manual")
	return Ok(true)
}

// UpdateIsrSetting
// @Summary update isr setting
// @Schemes
// @Description update isr setting
// @Tags Isr
// @Accept json
// @Produce json
// @Param        isr   body      domain.IsrSetting   true  "isr"
// @Success 200 {object} domain.Draft
// @Router /api/admin/isr [Put]
func (i *IsrRoute) UpdateIsrSetting(c *gin.Context) *R {
	isrSetting := &domain.IsrSetting{}
	if err := c.Bind(isrSetting); err != nil {
		return InternalError(err)
	}
	setting, err := i.SettingService.SaveOrUpdateIsrSetting(isrSetting)
	if err != nil {
		return InternalError(err)
	}
	return Ok(setting)
}

// GetIsrSetting
// @Summary get isr setting
// @Schemes
// @Description get isr setting
// @Tags Isr
// @Accept json
// @Produce json
// @Success 200 {object} domain.IsrSetting
// @Router /api/admin/isr [Get]
func (i *IsrRoute) GetIsrSetting(c *gin.Context) *R {
	isrSetting := i.SettingService.FindIsrSetting()
	return Ok(isrSetting)
}

func (i *IsrRoute) Register(r *gin.Engine) {
	r.POST(isrPathPrefix, Handle(i.ActiveAll))
	r.PUT(isrPathPrefix, Handle(i.UpdateIsrSetting))
	r.GET(isrPathPrefix, Handle(i.GetIsrSetting))
}
