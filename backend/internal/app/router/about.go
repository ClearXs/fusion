package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/svr"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

const AboutPathPrefix = "/api/admin/meta/about"

type AboutRouter struct {
	Cfg     *config.Config
	MetaSvr *svr.MetaService
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
func (router *AboutRouter) GetAbout(c *gin.Context) *R {
	meta := router.MetaSvr.GetMeta()
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
func (router *AboutRouter) UpdateAbout(c *gin.Context) *R {
	if router.Cfg.Demo {
		return InternalError(errors.New("演示站禁止修改此项！"))
	}
	var aboutCredential *credential.AboutCredential
	err := c.Bind(aboutCredential)
	if err != nil {
		return InternalError(err)
	}
	succeed := router.MetaSvr.UpdateAboutContent(aboutCredential.Content)
	return Ok(succeed)
}

func (router *AboutRouter) Register(r *gin.Engine) {
	r.GET(AboutPathPrefix, Handle(router.GetAbout))
	r.PUT(AboutPathPrefix, Handle(router.UpdateAbout))
}
