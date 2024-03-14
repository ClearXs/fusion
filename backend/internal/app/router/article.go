package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/web"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

const ArticlePathPrefix = "/api/admin/article"

type ArticleRouter struct {
	Cfg        *config.Config
	ArticleSvr *svr.ArticleService
}

var ArticleRouterSet = wire.NewSet(wire.Struct(new(ArticleRouter), "*"))

// GetArticleByOption
// @Summary 根据参数获取文章
// @Schemes
// @Description from option obtain articles
// @Tags Article
// @Accept json
// @Produce json
// @Param        page               query      int        false       "page"               -1
// @Param        pageSize           query      int        false       "pageSize"            5
// @Param        toListView         query      bool       false       "toListView"          false
// @Param        category           query      string     false       "category"
// @Param        tags               query      string     false       "page"
// @Param        title              query      int        false       "title"
// @Param        sortCreateAt       query      int        false       "sortCreateAt"
// @Param        sortTop            query      int        false       "sortTop"
// @Param        sortViewer         query      int        false       "sortViewer"
// @Param        startTime          query      int        false       "startTime"
// @Param        endTime            query      int        false       "endTime"
// @Success 200 {object} domain.ArticlePageResult
// @Router /api/admin/article [Get]
func (router *ArticleRouter) GetArticleByOption(c *gin.Context) *R {
	page := web.ParseNumberForQuery(c, "page", -1)
	pageSize := web.ParseNumberForQuery(c, "pageSize", 5)
	toListView := web.ParseBoolForQuery(c, "toListView", false)
	category := c.Query("category")
	tags := c.Query("tags")
	title := c.Query("title")
	sortCreateAt := c.Query("sortCreateAt")
	sortTop := c.Query("sortTop")
	sortViewer := c.Query("sortViewer")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	option := credential.ArticleSearchOptionCredential{
		Page:          page,
		PageSize:      pageSize,
		ToListView:    toListView,
		Category:      category,
		Tags:          tags,
		Title:         title,
		SortCreatedAt: sortCreateAt,
		SortTop:       sortTop,
		SortViewer:    sortViewer,
		StartTime:     startTime,
		EndTime:       endTime,
	}
	result := router.ArticleSvr.GetByOption(option, false)
	return Ok(result)
}

// GetOneById
// @Summary 根据id获取article
// @Schemes
// @Description 根据id获取article
// @Tags Article
// @Accept json
// @Produce json
// @Success 200 {object} domain.ArticlePageResult
// @Router /api/admin/article/:id [Get]
func (router *ArticleRouter) GetOneById(c *gin.Context) *R {
	id := web.ParseNumberForPath(c, "id", -1)
	article := router.ArticleSvr.GetById(int64(id))
	return Ok(article)
}

// UpdateArticleById
// @Summary 根据id更新文档
// @Schemes
// @Description 根据id获取article
// @Tags Article
// @Accept json
// @Produce json
// @Param       request     body      domain.Article   true     "query params"
// @Success 200 {object} bool
// @Router /api/admin/article/:id [PUT]
func (router *ArticleRouter) UpdateArticleById(c *gin.Context) *R {
	if router.Cfg.Demo {
		return Error(401, errors.New("演示站禁止修改文章！！"))
	}
	var article = &domain.Article{}
	err := c.Bind(article)
	if err != nil {
		return InternalError(err)
	}
	id := web.ParseNumberForPath(c, "id", -1)
	updated := router.ArticleSvr.UpdateById(int64(id), article)
	return Ok(updated)
}

// CreateArticle
// @Summary 创建文章
// @Schemes
// @Description 创建文章
// @Tags Article
// @Accept json
// @Produce json
// @Param       request     body      domain.Article   true     "query params"
// @Success 200 {object} bool
// @Router /api/admin/article/:id [POST]
func (router *ArticleRouter) CreateArticle(c *gin.Context) *R {
	if router.Cfg.Demo {
		return Error(401, errors.New("演示站禁止创建文章！！"))
	}
	var article = &domain.Article{}
	err := c.Bind(article)
	if err != nil {
		return InternalError(err)
	}
	create, err := router.ArticleSvr.Create(article)
	if err != nil {
		return InternalError(err)
	}
	return Ok(create)
}

// DeleteArticle
// @Summary 删除文章
// @Schemes
// @Description 删除文章
// @Tags Article
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/article/:id [DELETE]
func (router *ArticleRouter) DeleteArticle(c *gin.Context) *R {
	if router.Cfg.Demo {
		return Error(401, errors.New("演示站禁止删除文章！！！"))
	}
	id := web.ParseNumberForPath(c, "id", -1)
	deleted := router.ArticleSvr.DeleteById(int64(id))
	return Ok(deleted)
}

// GetArticlesByLink
// @Summary 根据link获取文章集合
// @Schemes
// @Description 根据link获取文章集合
// @Tags Article
// @Accept json
// @Produce json
// @Success 200 {object} domain.Article
// @Router /api/admin/article/searchByLink [POST]
func (router *ArticleRouter) GetArticlesByLink(c *gin.Context) *R {
	var searchArticleLink = &credential.ArticleSearchLinkCredential{}
	err := c.Bind(searchArticleLink)
	if err != nil {
		return InternalError(err)
	}
	articles := router.ArticleSvr.GetByLink(searchArticleLink.Link)
	return Ok(articles)
}

func (router *ArticleRouter) Register(r *gin.Engine) {
	r.GET(ArticlePathPrefix, Handle(router.GetArticleByOption))
	r.GET(ArticlePathPrefix+"/:id", Handle(router.GetOneById))
	r.PUT(ArticlePathPrefix+"/:id", Handle(router.UpdateArticleById))
	r.POST(ArticlePathPrefix+"", Handle(router.CreateArticle))
	r.DELETE(ArticlePathPrefix+"/:id", Handle(router.DeleteArticle))
	r.POST(ArticlePathPrefix+"/searchByLink", Handle(router.GetArticlesByLink))
}
