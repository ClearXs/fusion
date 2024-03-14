package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/svr"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type CustomPageRouter struct {
	Cfg           *config.Config
	StaticService *svr.StaticService
}

var CustomPageRouterSet = wire.NewSet(wire.Struct(new(CustomPageRouter), "*"))

func (c *CustomPageRouter) Register(r *gin.Engine) {

}
