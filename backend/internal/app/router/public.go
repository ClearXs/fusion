package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/env"
	"cc.allio/fusion/pkg/web"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/url"
)

const PublicPathPrefix = "/api/public"

type PublicRouter struct {
	Cfg               *config.Config
	ArticleService    *svr.ArticleService
	TagService        *svr.TagService
	MetaService       *svr.MetaService
	ViewerService     *svr.ViewerService
	VisitService      *svr.VisitService
	SettingService    *svr.SettingService
	CustomPageService *svr.CustomPageService
	CategoryService   *svr.CategoryService
}

var PublicRouterSet = wire.NewSet(wire.Struct(new(PublicRouter), "*"))

// getAll
// @Summary get custom page
// @Schemes
// @Description get custom page
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {object} []domain.CustomPage
// @Router /api/public/customPage/all [Get]
func (p *PublicRouter) getAll(c *gin.Context) *R {
	customPages := p.CustomPageService.GetAll()
	return Ok(customPages)
}

// getOneByPath
// @Summary get custom page
// @Schemes
// @Description get custom page
// @Tags Public
// @Accept json
// @Produce json
// @Param        path   path      string  true  "path"
// @Success 200 {object} domain.CustomPage
// @Router /api/public/customPage [Get]
func (p *PublicRouter) getOneByPath(c *gin.Context) *R {
	path := c.Query("path")
	customPage := p.CustomPageService.GetByPath(path)
	return Ok(customPage)
}

// getArticleByIdOrPathname
// @Summary get article by id or pathname
// @Schemes
// @Description get article by id or pathname
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {object} credential.AlternateArticle
// @Router /api/public/search [Get]
func (p *PublicRouter) getArticleByIdOrPathname(c *gin.Context) *R {
	idOrPathname := c.Param("id")
	alternateArticle, err := p.ArticleService.GetArticleByIdOrPathnameWithAlternate(idOrPathname)
	if err != nil {
		return InternalError(err)
	}
	return Ok(alternateArticle)
}

// getArticleByIdOrPathnameWithPassword
// @Summary get article by id or pathname
// @Schemes
// @Description get article by id or pathname
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {object} domain.Article
// @Router /api/public/article/:id [Post]
func (p *PublicRouter) getArticleByIdOrPathnameWithPassword(c *gin.Context) *R {
	idOrPathname := c.Param("id")
	articlePasswordCredential := &credential.ArticlePasswordCredential{}
	if err := c.Bind(articlePasswordCredential); err != nil {
		return InternalError(err)
	}
	article := p.ArticleService.GetArticleByIdOrPathnameWithPassword(idOrPathname, articlePasswordCredential.Password)
	return Ok(article)
}

// getArticleByIdOrPathnameWithPassword
// @Summary get article by id or pathname
// @Schemes
// @Description get article by id or pathname
// @Tags Public
// @Accept json
// @Produce json
// @Param        value   path      string  true  "value"
// @Success 200 {object} []domain.Article
// @Router /api/public/search [Get]
func (p *PublicRouter) searchArticle(c *gin.Context) *R {
	text := c.Query("value")
	articles := p.ArticleService.SearchByText(text, false)
	return Ok(articles)
}

// addViewer
// @Summary add system viewer
// @Schemes
// @Description add system viewer
// @Tags Public
// @Accept json
// @Produce json
// @Param        isNew         path      bool  true  "isNew"
// @Param        isNewByPath   path      bool  true  "isNewByPath"
// @Success 200 {object} []domain.Article
// @Router /api/public/viewer [Post]
func (p *PublicRouter) addViewer(c *gin.Context) *R {
	refer := c.Request.Header.Get("refer")
	uri, err := url.ParseRequestURI(refer)
	if err != nil {
		return InternalError(err)
	}
	path := uri.Path
	if path == "" {
		return InternalError(errors.New("not found 'refer' pathname"))
	}
	isNew := web.ParseBoolForQuery(c, "isNew", false)
	isNewByPath := web.ParseBoolForQuery(c, "isNewByPath", false)
	dataViewer := p.MetaService.AddViewer(isNew, path, isNewByPath)
	return Ok(dataViewer)
}

// getViewer
// @Summary get system viewer
// @Schemes
// @Description get system viewer
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {object} domain.DataViewer
// @Router /api/public/viewer [Get]
func (p *PublicRouter) getViewer(c *gin.Context) *R {
	viewer := p.MetaService.GetViewer()
	return Ok(viewer)
}

// getViewerByArticleIdOrPathname
// @Summary get Visit by article id or pathname
// @Schemes
// @Description get Visit by article id or pathname
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {object} domain.Visit
// @Router /api/public/article/viewer/:id [Get]
func (p *PublicRouter) getViewerByArticleIdOrPathname(c *gin.Context) *R {
	idOrPathname := c.Param("id")
	visit := p.VisitService.GetByArticleIdOrPathname(idOrPathname)
	return Ok(visit)
}

// getArticlesByTagName
// @Summary get articles by tag name
// @Schemes
// @Description get articles by tag name
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {object} []credential.PublicArticle
// @Router /api/public/tag/:name [Get]
func (p *PublicRouter) getArticlesByTagName(c *gin.Context) *R {
	tagName := c.Param("name")
	articles := p.TagService.GetArticlesByTagName(tagName, false)
	publicArticles := make([]*credential.PublicArticle, 0)
	for _, article := range articles {
		publicArticle := &credential.PublicArticle{
			Id:        article.Id,
			Title:     article.Title,
			Content:   article.Content,
			Tags:      article.Tags,
			Category:  article.Category,
			Top:       article.Top,
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		}
		publicArticles = append(publicArticles, publicArticle)
	}
	return Ok(publicArticles)
}

// getByOption
// @Summary 根据参数获取文章
// @Schemes
// @Description from option obtain articles
// @Tags Public
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
// @Router /api/public/article [Get]
func (p *PublicRouter) getByOption(c *gin.Context) *R {
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
	result := p.ArticleService.GetByOption(option, false)
	return Ok(result)
}

// getTimeLineInfo
// @Summary get timeline
// @Schemes
// @Description get timeline
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Router /api/public/timeline [Get]
func (p *PublicRouter) getTimeLineInfo(c *gin.Context) *R {
	timeLine := p.ArticleService.GetTimeLine()
	return Ok(timeLine)
}

// getArticlesByCategory
// @Summary get articles by category
// @Schemes
// @Description get articles by category
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Router /api/public/category [Get]
func (p *PublicRouter) getArticlesByCategory(c *gin.Context) *R {
	article := p.CategoryService.GetCategoriesWithArticle(false)
	return Ok(article)
}

// getArticlesByTag
// @Summary get articles by tag
// @Schemes
// @Description get articles by tag
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Router /api/public/tag [Get]
func (p *PublicRouter) getArticlesByTag(c *gin.Context) *R {
	article := p.TagService.GetTagsWithArticle(false)
	return Ok(article)
}

// getBuildMeta
// @Summary get site meta
// @Schemes
// @Description get site meta
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {object} credential.SiteMeta
// @Router /api/public/meta [Get]
func (p *PublicRouter) getBuildMeta(c *gin.Context) *R {
	tags := p.TagService.GetAllTags(false)
	meta := p.MetaService.GetMeta()
	categories := p.CategoryService.GetAllCategories()
	menus := p.SettingService.FindMenuSettings()
	totalArticles := p.ArticleService.GetTotalNum(false)
	totalWords := p.MetaService.GetTotalWords()
	layoutSetting := p.SettingService.FindLayoutSetting()
	version := env.Version
	siteMeta := &credential.SiteMeta{
		Version: version,
		Tags:    tags,
		Meta: struct {
			Id             int64                `json:"id"`
			Links          []*domain.LinkItem   `json:"links"`
			Socials        []*domain.SocialItem `json:"socials"`
			Menus          []*domain.MenuItem   `json:"menus"`
			Rewards        []*domain.RewardItem `json:"rewards"`
			About          *domain.About        `json:"about"`
			SiteInfo       *domain.SiteInfo     `json:"siteInfo"`
			Viewer         int64                `json:"viewer"`
			Visited        int64                `json:"visited"`
			TotalWordCount int64                `json:"totalWordCount"`
			Categories     []*domain.Category   `json:"categories"`
		}(struct {
			Id             int64
			Links          []*domain.LinkItem
			Socials        []*domain.SocialItem
			Menus          []*domain.MenuItem
			Rewards        []*domain.RewardItem
			About          *domain.About
			SiteInfo       *domain.SiteInfo
			Viewer         int64
			Visited        int64
			TotalWordCount int64
			Categories     []*domain.Category
		}{Id: meta.Id,
			Links:          meta.Links,
			Socials:        meta.Socials,
			Menus:          meta.Menus,
			Rewards:        meta.Rewards,
			About:          meta.About,
			SiteInfo:       meta.SiteInfo,
			Viewer:         meta.Viewer,
			Visited:        meta.Visited,
			TotalWordCount: meta.TotalWordCount,
			Categories:     categories}),
		Menus:         menus,
		TotalArticles: totalArticles,
		TotalWords:    totalWords,
		LayoutSetting: layoutSetting,
	}
	return Ok(siteMeta)
}

func (p *PublicRouter) Register(r *gin.Engine) {
	r.GET(PublicPathPrefix+"/customPage/all", Handle(p.getAll))
	r.GET(PublicPathPrefix+"/customPage", Handle(p.getOneByPath))
	r.GET(PublicPathPrefix+"/article/:id", Handle(p.getArticleByIdOrPathname))
	r.POST(PublicPathPrefix+"/article/:id", Handle(p.getArticleByIdOrPathnameWithPassword))
	r.GET(PublicPathPrefix+"/article/viewer/:id", Handle(p.getViewerByArticleIdOrPathname))
	r.GET(PublicPathPrefix+"/article", Handle(p.getByOption))

	r.GET(PublicPathPrefix+"/search", Handle(p.searchArticle))
	r.GET(PublicPathPrefix+"/viewer", Handle(p.getViewer))
	r.POST(PublicPathPrefix+"/viewer", Handle(p.addViewer))

	r.GET(PublicPathPrefix+"/tag/:name", Handle(p.getArticlesByTagName))
	r.GET(PublicPathPrefix+"/tag", Handle(p.getArticlesByTag))

	r.GET(PublicPathPrefix+"/timeline", Handle(p.getTimeLineInfo))

	r.GET(PublicPathPrefix+"/category", Handle(p.getArticlesByCategory))

	r.GET(PublicPathPrefix+"/meta", Handle(p.getBuildMeta))
}
