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

const SocialPathPrefix = "/api/admin/meta/social"

type SocialRoute struct {
	Cfg         *config.Config
	MetaService *svr.MetaService
	Isr         *event.IsrEventBus
}

var SocialRouterSet = wire.NewSet(wire.Struct(new(SocialRoute), "*"))

// GetSocials
// @Summary get socials
// @Schemes
// @Description get socials
// @Tags Social
// @Accept json
// @Produce json
// @Success 200 {object} []domain.SocialItem
// @Router /api/admin/meta/social [Get]
func (s *SocialRoute) GetSocials(c *gin.Context) *R {
	socials := s.MetaService.GetSocials()
	return Ok(socials)
}

// GetSocialTypes
// @Summary get socials
// @Schemes
// @Description get socials
// @Tags Social
// @Accept json
// @Produce json
// @Success 200 {object} []domain.SocialItem
// @Router /api/admin/meta/social/types [Get]
func (s *SocialRoute) GetSocialTypes(c *gin.Context) *R {
	socials := s.MetaService.GetDefaultSocials()
	return Ok(socials)
}

// UpdateSocial
// @Summary update social
// @Schemes
// @Description update social
// @Tags Social
// @Accept json
// @Produce json
// @Param        social   body      domain.SocialItem   true  "social"
// @Success 200 {object} bool
// @Router /api/admin/meta/social [Put]
func (s *SocialRoute) UpdateSocial(c *gin.Context) *R {
	if s.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	social := &domain.SocialItem{}
	if err := c.Bind(social); err != nil {
		return InternalError(err)
	}
	successed, err := s.MetaService.SaveOrUpdateSocial(social)
	if err != nil {
		return InternalError(err)
	}
	s.Isr.ActiveAll("trigger incremental rendering by social update")
	return Ok(successed)
}

// CreateSocial
// @Summary create social
// @Schemes
// @Description create social
// @Tags Social
// @Accept json
// @Produce json
// @Param        social   body      domain.SocialItem   true  "social"
// @Success 200 {object} bool
// @Router /api/admin/meta/social [Post]
func (s *SocialRoute) CreateSocial(c *gin.Context) *R {
	if s.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	social := &domain.SocialItem{}
	if err := c.Bind(social); err != nil {
		return InternalError(err)
	}
	successed, err := s.MetaService.SaveOrUpdateSocial(social)
	if err != nil {
		return InternalError(err)
	}
	s.Isr.ActiveAll("trigger incremental rendering by social create")
	return Ok(successed)
}

// DeleteSocial
// @Summary delete social by type name
// @Schemes
// @Description delete social by type name
// @Tags Social
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/meta/social/:type [Delete]
func (s *SocialRoute) DeleteSocial(c *gin.Context) *R {
	if s.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	typeName := c.Query("type")
	successed, err := s.MetaService.DeleteSocialByType(typeName)
	if err != nil {
		return InternalError(err)
	}
	s.Isr.ActiveAll("trigger incremental rendering by social delete")
	return Ok(successed)
}

func (s *SocialRoute) Register(r *gin.Engine) {
	r.GET(SocialPathPrefix, Handle(s.GetSocials))
	r.GET(SocialPathPrefix+"/types", Handle(s.GetSocialTypes))
	r.PUT(SocialPathPrefix, Handle(s.UpdateSocial))
	r.POST(SocialPathPrefix, Handle(s.CreateSocial))
	r.DELETE(SocialPathPrefix+"/:type", Handle(s.DeleteSocial))
}
