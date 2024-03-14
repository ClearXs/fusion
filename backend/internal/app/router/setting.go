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

const SettingPathPrefix = "/api/admin/setting"

type SettingRouter struct {
	Cfg            *config.Config
	SettingService *svr.SettingService
}

var SettingRouterSet = wire.NewSet(wire.Struct(new(SettingRouter), "*"))

// GetStaticSetting
// @Summary get static setting
// @Schemes
// @Description get static setting
// @Tags Setting
// @Accept json
// @Produce json
// @Success 200 {object} domain.StaticSetting
// @Router /api/admin/setting/static [Get]
func (s *SettingRouter) GetStaticSetting(c *gin.Context) *R {
	staticSetting := s.SettingService.FindStaticSetting()
	return Ok(staticSetting)
}

// UpdateStaticSetting
// @Summary save or update static setting
// @Schemes
// @Description save or update static setting
// @Tags Setting
// @Accept json
// @Produce json
// @Param        static   body      domain.StaticSetting   true  "static"
// @Success 200 {object} bool
// @Router /api/admin/setting/static [Put]
func (s *SettingRouter) UpdateStaticSetting(c *gin.Context) *R {
	if s.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	static := &domain.StaticSetting{}
	if err := c.Bind(static); err != nil {
		return InternalError(err)
	}
	successed, err := s.SettingService.SaveOrUpdateStaticSetting(static)
	if err != nil {
		return InternalError(err)
	}
	return Ok(successed)
}

// GetWalineSetting
// @Summary get waline setting
// @Schemes
// @Description get waline setting
// @Tags Setting
// @Accept json
// @Produce json
// @Success 200 {object} domain.WalineSetting
// @Router /api/admin/setting/waline [Get]
func (s *SettingRouter) GetWalineSetting(c *gin.Context) *R {
	waline := s.SettingService.FindWalineSetting()
	return Ok(waline)
}

// UpdateWalineSetting
// @Summary save or update waline setting
// @Schemes
// @Description save or update waline setting
// @Tags Setting
// @Accept json
// @Produce json
// @Param        waline   body      domain.WalineSetting   true  "waline"
// @Success 200 {object} bool
// @Router /api/admin/setting/waline [Put]
func (s *SettingRouter) UpdateWalineSetting(c *gin.Context) *R {
	if s.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	waline := &domain.WalineSetting{}
	if err := c.Bind(waline); err != nil {
		return InternalError(err)
	}
	successed, err := s.SettingService.SaveOrUpdateWalineSetting(waline)
	if err != nil {
		return InternalError(err)
	}
	return Ok(successed)
}

// GetLayoutSetting
// @Summary get layout setting
// @Schemes
// @Description get layout setting
// @Tags Setting
// @Accept json
// @Produce json
// @Success 200 {object} domain.LayoutSetting
// @Router /api/admin/setting/layout [Get]
func (s *SettingRouter) GetLayoutSetting(c *gin.Context) *R {
	layout := s.SettingService.FindLayoutSetting()
	return Ok(layout)
}

// UpdateLayoutSetting
// @Summary save or update layout setting
// @Schemes
// @Description save or update layout setting
// @Tags Setting
// @Accept json
// @Produce json
// @Param        layout   body      domain.LayoutSetting   true  "layout"
// @Success 200 {object} bool
// @Router /api/admin/setting/layout [Put]
func (s *SettingRouter) UpdateLayoutSetting(c *gin.Context) *R {
	if s.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	layout := &domain.LayoutSetting{}
	if err := c.Bind(layout); err != nil {
		return InternalError(err)
	}
	successed, err := s.SettingService.SaveOrUpdateLayoutSetting(layout)
	if err != nil {
		return InternalError(err)
	}
	return Ok(successed)
}

// GetLoginSetting
// @Summary get login setting
// @Schemes
// @Description get login setting
// @Tags Setting
// @Accept json
// @Produce json
// @Success 200 {object} domain.LoginSetting
// @Router /api/admin/setting/login [Get]
func (s *SettingRouter) GetLoginSetting(c *gin.Context) *R {
	login := s.SettingService.FindLoginSetting()
	return Ok(login)
}

// UpdateLoginSetting
// @Summary save or update login setting
// @Schemes
// @Description save or update login setting
// @Tags Setting
// @Accept json
// @Produce json
// @Param        login   body      domain.LoginSetting   true  "login"
// @Success 200 {object} bool
// @Router /api/admin/setting/login [Put]
func (s *SettingRouter) UpdateLoginSetting(c *gin.Context) *R {
	if s.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	login := &domain.LoginSetting{}
	if err := c.Bind(login); err != nil {
		return InternalError(err)
	}
	successed, err := s.SettingService.SaveOrUpdateLoginSetting(login)
	if err != nil {
		return InternalError(err)
	}
	return Ok(successed)
}

func (s *SettingRouter) Register(r *gin.Engine) {
	r.GET(SettingPathPrefix+"/static", Handle(s.GetStaticSetting))
	r.PUT(SettingPathPrefix+"/static", Handle(s.UpdateStaticSetting))

	r.GET(SettingPathPrefix+"/waline", Handle(s.GetWalineSetting))
	r.PUT(SettingPathPrefix+"/waline", Handle(s.UpdateWalineSetting))

	r.GET(SettingPathPrefix+"/layout", Handle(s.GetLayoutSetting))
	r.PUT(SettingPathPrefix+"/layout", Handle(s.UpdateLayoutSetting))

	r.GET(SettingPathPrefix+"/login", Handle(s.GetLoginSetting))
	r.PUT(SettingPathPrefix+"/login", Handle(s.UpdateLoginSetting))
}
