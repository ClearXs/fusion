package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/ip"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

const CaddyPathPrefix = "/api/admin/caddy"

type CaddyRouter struct {
	Cfg             *config.Config
	SettingsService *svr.SettingService
	CaddyService    *svr.CaddyService
}

var CaddyRouterSet = wire.NewSet(wire.Struct(new(CaddyRouter), "*"))

// GetHttpsSettings
// @Summary acquire https settings
// @Schemes
// @Description acquire https settings
// @Tags Caddy
// @Accept json
// @Produce json
// @Success 200 {object} domain.HttpSetting
// @Router /api/admin/caddy/https [Get]
func (caddy *CaddyRouter) GetHttpsSettings(c *gin.Context) *R {
	httpsSetting := caddy.SettingsService.FindHttpsSetting()
	return Ok(httpsSetting)
}

// updateHttpsSettings
// @Summary update https settings
// @Schemes
// @Description update https settings
// @Tags Caddy
// @Accept json
// @Produce json
// @Param        httpsSettingCredential   body      credential.HttpsSettingCredential   true  "httpsSettingCredential"
// @Success 200 {object} bool
// @Router /api/admin/caddy/https [Put]
func (caddy *CaddyRouter) updateHttpsSettings(c *gin.Context) *R {
	if caddy.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	httpsSettingCredential := &credential.HttpsSettingCredential{}
	if err := c.Bind(httpsSettingCredential); err != nil {
		return InternalError(err)
	}
	successed := caddy.CaddyService.SetRedirect(httpsSettingCredential.Redirect)
	if !successed {
		return InternalError(errors.New("set redirect failed"))
	}
	successed, err := caddy.SettingsService.SaveOrUpdateHttpsSetting(httpsSettingCredential)
	if err != nil {
		return InternalError(err)
	}
	return Ok(successed)
}

// askOnDemand
// @Summary acquire https settings
// @Schemes
// @Description acquire https settings
// @Tags Caddy
// @Accept json
// @Produce json
// @Param        domain   path    string   true  "domain"
// @Success 200 {object} nil
// @Router /api/admin/caddy/ask [Get]
func (caddy *CaddyRouter) askOnDemand(c *gin.Context) *R {
	domain := c.Param("domain")
	isIpV4 := ip.IsIpV4(domain)
	if !isIpV4 {
		return OkMessage(nil, "is domain, on demand https")
	}
	return InternalError(errors.New("error domain"))
}

// getLog
// @Summary get caddy log
// @Schemes
// @Description get caddy log
// @Tags Caddy
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /api/admin/caddy/log [Get]
func (caddy *CaddyRouter) getLog(c *gin.Context) *R {
	// TODO 未知/var/log/caddy.log是什么用
	log, err := caddy.CaddyService.GetLog()
	if err != nil {
		return InternalError(err)
	}
	return Ok(log)
}

// clearLog
// @Summary clear caddy log
// @Schemes
// @Description clear caddy log
// @Tags Caddy
// @Accept json
// @Produce json
// @Success 200 {object} nil
// @Router /api/admin/caddy/log [Delete]
func (caddy *CaddyRouter) clearLog(c *gin.Context) *R {
	// TODO 未知/var/log/caddy.log是什么用
	if caddy.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	caddy.CaddyService.ClearLog()
	return OkMessage(nil, "clear caddy log success")
}

// getCaddyConfig
// @Summary get caddy config
// @Schemes
// @Description get caddy config
// @Tags Caddy
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /api/admin/caddy/config [Get]
func (caddy *CaddyRouter) getCaddyConfig(c *gin.Context) *R {
	// TODO 未知 http://127.0.0.1:2019/config是什么接口
	cfg, err := caddy.CaddyService.GetConfig()
	if err != nil {
		return InternalError(err)
	}
	return Ok(cfg)
}

func (caddy *CaddyRouter) Register(r *gin.Engine) {
	r.GET(CaddyPathPrefix+"https", Handle(caddy.GetHttpsSettings))
	r.PUT(CaddyPathPrefix+"https", Handle(caddy.updateHttpsSettings))

	r.GET(CaddyPathPrefix+"ask", Handle(caddy.askOnDemand))

	r.GET(CaddyPathPrefix+"log", Handle(caddy.getLog))
	r.DELETE(CaddyPathPrefix+"log", Handle(caddy.clearLog))

	r.GET(CaddyPathPrefix+"config", Handle(caddy.getCaddyConfig))
}
