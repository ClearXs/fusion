package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/env"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"time"
)

const MetaPathPrefix = "/api/admin/meta"

type MetaRouter struct {
	Cfg         *config.Config
	MetaService *svr.MetaService
}

var MetaRouterSet = wire.NewSet(wire.Struct(new(MetaRouter), "*"))

// GetAllMeta
// @Summary get meta info
// @Schemes
// @Description get meta info
// @Tags Meta
// @Accept json
// @Produce json
// @Success 200 {object} credential.MetaCredential
// @Router /api/admin/meta [Get]
func (m *MetaRouter) GetAllMeta(c *gin.Context) *R {
	meta := m.MetaService.GetMeta()
	return Ok(&credential.MetaCredential{
		Version:       env.Version,
		LatestVersion: env.Version,
		UpdateAt:      time.Now(),
		BaseUrl:       meta.SiteInfo.BaseUrl,
		EnableComment: meta.SiteInfo.EnableComment,
		AllowDomains:  env.FusionAllowDomains,
	})
}

func (m *MetaRouter) Register(r *gin.Engine) {
	r.GET(MetaPathPrefix, Handle(m.GetAllMeta))
}
