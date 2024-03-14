package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/cache"
	"cc.allio/fusion/pkg/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
	"net/http"
	"time"
)

const AuthPathPrefix = "/api/admin/auth"

type AuthRouter struct {
	Cfg      *config.Config
	AuthSvr  *svr.AuthService
	UserSvr  *svr.UserService
	TokenSvr *svr.TokenService
}

var AuthRouterSet = wire.NewSet(wire.Struct(new(AuthRouter), "*"))

// Login
// @Summary 登陆
// @Schemes
// @Description login by username and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param        username   path      string  true  "username"
// @Param        password   path      string  true  "password"
// @Success 200 {object} svr.TokenUser
// @Router /api/admin/auth/login [Post]
func (router *AuthRouter) Login(c *gin.Context) *R {
	ip := utils.GetIpByRequest(c.Request)
	slog.Info("ip %s", ip)
	var loginCredential credential.LoginCredential
	err := c.ShouldBind(&loginCredential)
	if err != nil {
		return AuthenticationError(err)
	}
	tokenUser, err := router.AuthSvr.Login(loginCredential.Username, loginCredential.Password)
	if err != nil {
		return AuthenticationError(err)
	}
	return Ok(tokenUser)
}

// Logout
// @Summary 登出
// @Description logout
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {bool} bool
// @Router /api/admin/auth/logout [GET]
func (router *AuthRouter) Logout(c *gin.Context) *R {
	token := c.GetHeader("token")
	if lo.IsEmpty(token) {
		return AuthenticationError(errors.New("无登录凭证"))
	}
	succeed, err := router.TokenSvr.DisabledToken(token)
	if err != nil {
		return InternalError(err)
	}
	if !succeed {
		return InternalError(errors.New("登出失败"))
	} else {
		return OkMessage(succeed, "登出成功")
	}
}

// Restore
// @Summary 恢复密钥
// @Description Restore
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {bool} bool
// @Router /api/admin/auth/restore [POST]
func (router *AuthRouter) Restore(c *gin.Context) *R {
	var restoreCredential credential.RestoreCredential
	err := c.ShouldBind(&restoreCredential)
	if err != nil {
		return InternalError(err)
	}
	token := restoreCredential.Key
	keyInCache := cache.Get("restoreKey")
	if lo.IsEmpty(token) || token != keyInCache {
		return AuthenticationError(errors.New("恢复密钥错误"))
	}
	success, err := router.UserSvr.UpdateUser(&domain.UpdateUser{Name: restoreCredential.Username, Password: restoreCredential.Password})
	if err != nil || !success {
		return InternalError(err)
	}
	// 在前端清理 localStore 之后
	time.AfterFunc(60*time.Second, func() {
		router.TokenSvr.DisabledAll()
	})
	return OkMessage(nil, "重置成功!")
}

// UpdateUser
// @Summary 更新用户
// @Description 更新用户
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {bool} bool
// @Router /api/admin/auth [PUT]
func (router *AuthRouter) UpdateUser(c *gin.Context) *R {
	if router.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改账号密码！"))
	}
	var updateUser = &domain.UpdateUser{}
	if err := c.Bind(updateUser); err != nil {
		return InternalError(err)
	}
	updated, err := router.UserSvr.UpdateUser(updateUser)
	if err != nil {
		return InternalError(err)
	}
	return Ok(updated)
}

func (router *AuthRouter) Register(r *gin.Engine) {
	r.POST(AuthPathPrefix+"/login", Handle(router.Login))
	r.GET(AuthPathPrefix+"/logout", Handle(router.Logout))
	r.POST(AuthPathPrefix+"/restore", Handle(router.Restore))
	r.PUT(AuthPathPrefix+"", Handle(router.UpdateUser))
}
