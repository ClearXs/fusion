package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/svr"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

const UserPrefix = "/api/admin/user"

type UserRoute struct {
	Cfg     *config.Config
	UserSvr *svr.UserService
}

var UserRouterSet = wire.NewSet(wire.Struct(new(UserRoute), "*"))

// GetUserList
// @Summary 获取用户列表
// @Schemes
// @Description obtain user list
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} domain.User
// @Router /api/admin/user/list [Get]
func (u *UserRoute) GetUserList(c *gin.Context) *R {
	userList := u.UserSvr.GetUserList()
	return Ok(userList)
}

func (u *UserRoute) Register(r *gin.Engine) {
	r.GET(UserPrefix+"/list", Handle(u.GetUserList))
}
