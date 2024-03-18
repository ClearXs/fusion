package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/event"
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
	Isr        *event.IsrEventBus
	Script     *event.ScriptEngine
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
func (a *ArticleRouter) GetArticleByOption(c *gin.Context) *R {
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
	result := a.ArticleSvr.GetByOption(option, false)
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
func (a *ArticleRouter) GetOneById(c *gin.Context) *R {
	id := web.ParseNumberForPath(c, "id", -1)
	article := a.ArticleSvr.GetById(int64(id))
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
func (a *ArticleRouter) UpdateArticleById(c *gin.Context) *R {
	if a.Cfg.Demo {
		return Error(401, errors.New("演示站禁止修改文章！！"))
	}
	var article = &domain.Article{}
	err := c.Bind(article)
	if err != nil {
		return InternalError(err)
	}
	id := web.ParseNumberForPath(c, "id", -1)
	a.Script.DispatchBeforeUpdateArticleEvent(article)
	updated := a.ArticleSvr.UpdateById(int64(id), article)
	a.Isr.ActiveAll("trigger incremental rendering by update article", article)
	a.Script.DispatchAfterUpdateArticleEvent(article, updated)
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
func (a *ArticleRouter) CreateArticle(c *gin.Context) *R {
	if a.Cfg.Demo {
		return Error(401, errors.New("演示站禁止创建文章！！"))
	}
	var article = &domain.Article{}
	err := c.Bind(article)
	if err != nil {
		return InternalError(err)
	}
	a.Script.DispatchBeforeUpdateArticleEvent(article)
	create, err := a.ArticleSvr.Create(article)
	if err != nil {
		return InternalError(err)
	}
	a.Script.DispatchAfterUpdateArticleEvent(article, create)
	a.Isr.ActiveAll("trigger incremental rendering by create article", article)
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
func (a *ArticleRouter) DeleteArticle(c *gin.Context) *R {
	if a.Cfg.Demo {
		return Error(401, errors.New("演示站禁止删除文章！！！"))
	}
	id := web.ParseNumberForPath(c, "id", -1)
	deleted := a.ArticleSvr.DeleteById(int64(id))
	a.Isr.ActiveAll("trigger incremental rendering by delete article", id)
	a.Script.DispatchDeleteArticleEvent(id, deleted)
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
func (a *ArticleRouter) GetArticlesByLink(c *gin.Context) *R {
	var searchArticleLink = &credential.ArticleSearchLinkCredential{}
	err := c.Bind(searchArticleLink)
	if err != nil {
		return InternalError(err)
	}
	articles := a.ArticleSvr.GetByLink(searchArticleLink.Link)
	return Ok(articles)
}

func (a *ArticleRouter) Register(r *gin.Engine) {
	r.GET(ArticlePathPrefix, Handle(a.GetArticleByOption))
	r.GET(ArticlePathPrefix+"/:id", Handle(a.GetOneById))
	r.PUT(ArticlePathPrefix+"/:id", Handle(a.UpdateArticleById))
	r.POST(ArticlePathPrefix+"", Handle(a.CreateArticle))
	r.DELETE(ArticlePathPrefix+"/:id", Handle(a.DeleteArticle))
	r.POST(ArticlePathPrefix+"/searchByLink", Handle(a.GetArticlesByLink))
}
