package router

import (
	"cc.allio/fusion/docs"
	"cc.allio/fusion/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type Register interface {
	Register(r *gin.Engine)
}

type Router struct {
	AboutRouter        *AboutRouter
	AnalysisRouter     *AnalysisRouter
	ArticleRouter      *ArticleRouter
	AuthRouter         *AuthRouter
	BackupRouter       *BackupRouter
	UserRouter         *UserRouter
	CaddyRouter        *CaddyRouter
	CategoryRouter     *CategoryRouter
	CollaboratorRouter *CollaboratorRouter
	CustomPageRouter   *CustomPageRouter
	DraftRouter        *DraftRouter
	LinkRouter         *LinkRouter
	MenuRouter         *MenuRouter
	MetaRouter         *MetaRouter
	RewardRouter       *RewardRouter
	SettingRouter      *SettingRouter
	SiteRouter         *SiteRouter
	SocialRouter       *SocialRouter
	TagRouter          *TagRouter
	TokenRouter        *TokenRouter
	PublicRouter       *PublicRouter
	PipelineRouter     *PipelineRouter
	IsrRouter          *IsrRouter
	ImgRouter          *ImgRouter
}

var Set = wire.NewSet(
	AboutRouterSet,
	AnalysisRouterSet,
	ArticleRouterSet,
	AuthRouterSet,
	BackupRouterSet,
	UserRouterSet,
	CaddyRouterSet,
	CategoryRouterSet,
	CollaboratorRouterSet,
	CustomPageRouterSet,
	DraftRouterSet,
	LinkRouterSet,
	MenuRouterSet,
	MetaRouterSet,
	RewardRouterSet,
	SettingRouterSet,
	SiteRouterSet,
	SocialRouterSet,
	TagRouterSet,
	TokenRouterSet,
	PublicRouterSet,
	PipelineRouterSet,
	IsrRouterSet,
	ImgRouterSet,
	wire.Struct(new(Router), "*"),
)

// Init initiation system router
func (router *Router) Init(r *gin.Engine) {

	// middleware setup
	r.Use(gin.Recovery())
	r.Use(middleware.Logging())

	// swagger
	docs.SwaggerInfo.Title = "Fusion"
	docs.SwaggerInfo.Description = "This is a Fusion system."
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.Host = "localhost:5600"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// register system api
	{
		router.AboutRouter.Register(r)
		router.AnalysisRouter.Register(r)
		router.AuthRouter.Register(r)
		router.ArticleRouter.Register(r)
		router.BackupRouter.Register(r)
		router.UserRouter.Register(r)
		router.CaddyRouter.Register(r)
		router.CategoryRouter.Register(r)
		router.CollaboratorRouter.Register(r)
		router.CustomPageRouter.Register(r)
		router.DraftRouter.Register(r)
		router.LinkRouter.Register(r)
		router.MenuRouter.Register(r)
		router.MetaRouter.Register(r)
		router.RewardRouter.Register(r)
		router.SettingRouter.Register(r)
		router.SiteRouter.Register(r)
		router.SocialRouter.Register(r)
		router.TagRouter.Register(r)
		router.TokenRouter.Register(r)
		router.PublicRouter.Register(r)
		router.PipelineRouter.Register(r)
		router.IsrRouter.Register(r)
		router.ImgRouter.Register(r)
	}

	// route or method not found
	r.NoRoute(routeOrMethodNotFound)
	r.NoMethod(routeOrMethodNotFound)
}

func routeOrMethodNotFound(c *gin.Context) {
	c.String(http.StatusNotFound, "route not found")
}
