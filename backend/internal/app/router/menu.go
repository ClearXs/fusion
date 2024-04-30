package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/event"
	"cc.allio/fusion/internal/svr"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

const MenuPathPrefix = "/api/admin/meta/menu"

type MenuRoute struct {
	Cfg             *config.Config
	SettingsService *svr.SettingService
	Isr             *event.IsrEventBus
}

var MenuRouterSet = wire.NewSet(wire.Struct(new(MenuRoute), "*"))

// GetMenu
// @Summary get menu by settings
// @Schemes
// @Description get menu by settings
// @Tags Menu
// @Accept json
// @Produce json
// @Success 200 {object} []domain.MenuItem
// @Router /api/admin/meta/menu [Get]
func (m *MenuRoute) GetMenu(c *gin.Context) *R {
	settings, err := m.SettingsService.FindMenuSettings()
	if err != nil {
		return InternalError(err)
	}
	return Ok(settings)
}

// UpdateMenu
// @Summary update menu
// @Schemes
// @Description update menu
// @Tags Menu
// @Accept json
// @Produce json
// @Param        menu   body      domain.MenuItem   true  "menu"
// @Success 200 {object} bool
// @Router /api/admin/meta/menu [Put]
func (m *MenuRoute) UpdateMenu(c *gin.Context) *R {
	if m.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	var menus []*domain.MenuItem
	if err := c.Bind(&menus); err != nil {
		return InternalError(err)
	}
	successd, err := m.SettingsService.SaveOrUpdateMenuSettings(menus)
	if err != nil {
		return InternalError(err)
	}
	m.Isr.ActiveAll("trigger incremental rendering by menu update")
	return Ok(successd)
}

func (m *MenuRoute) Register(r *gin.Engine) {
	r.GET(MenuPathPrefix, Handle(m.GetMenu))
	r.PUT(MenuPathPrefix, Handle(m.UpdateMenu))
}
