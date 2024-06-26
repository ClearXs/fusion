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

const SitePathPrefix = "/api/admin/meta/site"

type SiteRoute struct {
	Cfg         *config.Config
	MetaService *svr.MetaService
	Isr         *event.IsrEventBus
	Script      *event.ScriptEngine
}

var SiteRouterSet = wire.NewSet(wire.Struct(new(SiteRoute), "*"))

// GetSite
// @Summary get site
// @Schemes
// @Description get site
// @Tags Site
// @Accept json
// @Produce json
// @Success 200 {object} domain.SiteInfo
// @Router /api/admin/meta/site [Get]
func (s *SiteRoute) GetSite(c *gin.Context) *R {
	siteInfo := s.MetaService.GetSiteInfo()
	return Ok(siteInfo)
}

// UpdateSite
// @Summary update site info
// @Schemes
// @Description update site info
// @Tags Site
// @Accept json
// @Produce json
// @Param        site   body      domain.SiteInfo   true  "site"
// @Success 200 {object} bool
// @Router /api/admin/meta/site [Put]
func (s *SiteRoute) UpdateSite(c *gin.Context) *R {
	if s.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	siteInfo := &domain.SiteInfo{}
	if err := c.Bind(siteInfo); err != nil {
		return InternalError(err)
	}
	successed, err := s.MetaService.UpdateSiteInfo(siteInfo)
	if err != nil {
		return InternalError(err)
	}
	s.Isr.ActiveAll("trigger incremental rendering by update site info")
	s.Script.DispatchUpdateSiteInfoEvent(siteInfo)
	return Ok(successed)
}

func (s *SiteRoute) Register(r *gin.Engine) {
	r.GET(SitePathPrefix, Handle(s.GetSite))
	r.PUT(SitePathPrefix, Handle(s.UpdateSite))
}
