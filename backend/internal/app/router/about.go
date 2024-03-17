package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/event"
	"cc.allio/fusion/internal/svr"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

const AboutPathPrefix = "/api/admin/meta/about"

type AboutRouter struct {
	Cfg     *config.Config
	MetaSvr *svr.MetaService
	Isr     *event.IsrEventBus
}

var AboutRouterSet = wire.NewSet(wire.Struct(new(AboutRouter), "*"))

// GetAbout
// @Summary 获取About
// @Schemes
// @Description obtain about information
// @Tags Meta
// @Accept json
// @Produce json
// @Success 200 {object} domain.Meta
// @Router /api/admin/meta/about [Get]
func (a *AboutRouter) GetAbout(c *gin.Context) *R {
	meta := a.MetaSvr.GetMeta()
	return Ok(meta)
}

// UpdateAbout
// @Summary 更新 about
// @Schemes
// @Description obtain about information
// @Tags Meta
// @Accept json
// @Produce json
// @Param        content   path      string  true  "content"
// @Success 200 {object} bool
// @Router /api/admin/meta/about [Put]
func (a *AboutRouter) UpdateAbout(c *gin.Context) *R {
	if a.Cfg.Demo {
		return InternalError(errors.New("演示站禁止修改此项！"))
	}
	var aboutCredential *credential.AboutCredential
	err := c.Bind(aboutCredential)
	if err != nil {
		return InternalError(err)
	}
	succeed := a.MetaSvr.UpdateAboutContent(aboutCredential.Content)

	a.Isr.ActiveAbout("trigger about incremental by update about")
	return Ok(succeed)
}

func (a *AboutRouter) Register(r *gin.Engine) {
	r.GET(AboutPathPrefix, Handle(a.GetAbout))
	r.PUT(AboutPathPrefix, Handle(a.UpdateAbout))
}
