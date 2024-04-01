package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/svr"
	token2 "cc.allio/fusion/internal/token"
	"cc.allio/fusion/pkg/env"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"time"
)

const MetaPathPrefix = "/api/admin/meta"

type MetaRoute struct {
	Cfg         *config.Config
	MetaService *svr.MetaService
}

var MetaRouterSet = wire.NewSet(wire.Struct(new(MetaRoute), "*"))

// GetAllMeta
// @Summary get meta info
// @Schemes
// @Description get meta info
// @Tags Meta
// @Accept json
// @Produce json
// @Success 200 {object} credential.MetaCredential
// @Router /api/admin/meta [Get]
func (m *MetaRoute) GetAllMeta(c *gin.Context) *R {
	meta := m.MetaService.GetMeta()
	metaCredential := &credential.MetaCredential{
		Version:       env.Version,
		LatestVersion: env.Version,
		UpdateAt:      time.Now(),
		BaseUrl:       meta.SiteInfo.BaseUrl,
		EnableComment: meta.SiteInfo.EnableComment,
		AllowDomains:  env.FusionAllowDomains,
	}
	// find user by token string
	claims, err := token2.GetClaimsByRequest(c, m.Cfg.Token.SignedKey)
	if err != nil {
		return InternalError(err)
	}
	if tokenClaims, ok := claims.(*token2.Claims); ok {
		user := make(map[string]interface{})
		user["id"] = tokenClaims.Id
		user["name"] = tokenClaims.Name
		metaCredential.User = user
	}
	return Ok(metaCredential)
}

func (m *MetaRoute) Register(r *gin.Engine) {
	r.GET(MetaPathPrefix, Handle(m.GetAllMeta))
}
