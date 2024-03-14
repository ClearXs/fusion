package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/svr"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

const SitePathPrefix = "/api/admin/meta/site"

type SiteRouter struct {
	Cfg         *config.Config
	MetaService *svr.MetaService
}

var SiteRouterSet = wire.NewSet(wire.Struct(new(SiteRouter), "*"))

// GetSite
// @Summary get site
// @Schemes
// @Description get site
// @Tags Site
// @Accept json
// @Produce json
// @Success 200 {object} domain.SiteInfo
// @Router /api/admin/meta/site [Get]
func (s *SiteRouter) GetSite(c *gin.Context) *R {
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
func (s *SiteRouter) UpdateSite(c *gin.Context) *R {
	if s.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	siteInfo := &domain.SiteInfo{}
	successed, err := s.MetaService.UpdateSiteInfo(siteInfo)
	if err != nil {
		return InternalError(err)
	}
	return Ok(successed)
}

func (s *SiteRouter) Register(r *gin.Engine) {
	r.GET(SitePathPrefix, Handle(s.GetSite))
	r.PUT(SitePathPrefix, Handle(s.UpdateSite))
}
