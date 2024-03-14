package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/svr"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

const UserPrefix = "/api/admin/user"

type UserRouter struct {
	Cfg     *config.Config
	UserSvr *svr.UserService
}

var UserRouterSet = wire.NewSet(wire.Struct(new(UserRouter), "*"))

// GetUserList
// @Summary 获取用户列表
// @Schemes
// @Description obtain user list
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} domain.User
// @Router /api/admin/user/list [Get]
func (router *UserRouter) GetUserList(c *gin.Context) *R {
	userList := router.UserSvr.GetUserList()
	return Ok(userList)
}

func (router *UserRouter) Register(r *gin.Engine) {
	r.GET(UserPrefix+"/list", Handle(router.GetUserList))
}
