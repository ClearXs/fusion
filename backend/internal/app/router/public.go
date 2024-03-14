package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/svr"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

const PublicPathPrefix = "/api/public"

type PublicRouter struct {
	Cfg               *config.Config
	ArticleService    *svr.ArticleService
	TagService        *svr.TagService
	MetaService       *svr.MetaService
	VisitService      *svr.VisitService
	SettingService    *svr.SettingService
	CustomPageService *svr.CustomPageService
}

var PublicRouterSet = wire.NewSet(wire.Struct(new(PublicRouter), "*"))

func (p *PublicRouter) getAll(c *gin.Context) *R {
	customPages := p.CustomPageService.GetAll()
	return Ok(customPages)
}

func (p *PublicRouter) getOneByPath(c *gin.Context) *R {
	path := c.Query("path")
	customPage := p.CustomPageService.GetByPath(path)
	return Ok(customPage)
}

func (p *PublicRouter) getArticleByIdOrPathname(c *gin.Context) *R {
	idOrPathname := c.Param("id")
	alternateArticle, err := p.ArticleService.GetArticleByIdOrPathnameWithAlternate(idOrPathname)
	if err != nil {
		return InternalError(err)
	}
	return Ok(alternateArticle)
}

func (p *PublicRouter) getArticleByIdOrPathnameWithPassword(c *gin.Context) *R {
	idOrpathname := c.Param("id")
	articlePasswordCredential := &credential.ArticlePasswordCredential{}
	if err := c.Bind(articlePasswordCredential); err != nil {
		return InternalError(err)
	}
	article := p.ArticleService.GetArticleByIdOrPathnameWithPassword(idOrpathname, articlePasswordCredential.Password)
	return Ok(article)
}

func (p *PublicRouter) searchArticle(c *gin.Context) *R {

}

func (p *PublicRouter) addViewer(c *gin.Context) *R {

}

func (p *PublicRouter) getViewer(c *gin.Context) *R {

}

func (p *PublicRouter) getViewerByArticleIdOrPathname(c *gin.Context) *R {

}

func (p *PublicRouter) getArticlesByTagName(c *gin.Context) *R {

}

func (p *PublicRouter) getByOption(c *gin.Context) *R {

}

func (p *PublicRouter) getTimeLineInfo(c *gin.Context) *R {

}

func (p *PublicRouter) getArticlesByCategory(c *gin.Context) *R {

}

func (p *PublicRouter) getArticlesByTag(c *gin.Context) *R {

}

func (p *PublicRouter) getBuildMeta(c *gin.Context) *R {

}

func (p *PublicRouter) Register(r *gin.Engine) {

}
