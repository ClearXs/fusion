package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"github.com/google/wire"
	"time"
)

type AnalysisService struct {
	Cfg         *config.Config
	MetaSvr     *MetaService
	ArticleSvr  *ArticleService
	ViewerSvr   *ViewerService
	VisitSvr    *VisitService
	TagSvr      *TagService
	CategorySvr *CategoryService
}

var AnalysisServiceSet = wire.NewSet(wire.Struct(new(AnalysisService), "*"))

func (a *AnalysisService) GetWelcomePageData(tab string, overviewDataNum int64, viewerDataNum int64, articleTabDataNum int64) interface{} {
	if tab == "overview" {
		return a.GetOverViewTabData(overviewDataNum)
	}
	if tab == "viewer" {
		return a.GetViewerTabData(viewerDataNum)
	}
	if tab == "article" {
		return a.GetArticleTabData(articleTabDataNum)
	}
	return make(map[string]any)
}

func (a *AnalysisService) GetOverViewTabData(overviewDataNum int64) *domain.OverviewTabData {
	wordCount := a.MetaSvr.GetTotalWords()
	articleNum := a.ArticleSvr.GetTotalNum(true)
	viewerGrid := a.ViewerSvr.GetViewerGrid(overviewDataNum)
	siteInfo := a.MetaSvr.GetSiteInfo()
	return &domain.OverviewTabData{
		Total: struct {
			WordCount  int64 `json:"wordCount"`
			ArticleNum int64 `json:"articleNum"`
		}{WordCount: wordCount, ArticleNum: articleNum},
		Viewer: *viewerGrid,
		Link: struct {
			BaseUrl       string             `json:"baseUrl"`
			EnableComment domain.TrueOrFalse `json:"enableComment"`
		}{BaseUrl: siteInfo.BaseUrl, EnableComment: siteInfo.EnableComment},
	}
}

func (a *AnalysisService) GetViewerTabData(viewerDataNum int64) *domain.ViewerTabData {
	siteInfo := a.MetaSvr.GetSiteInfo()
	enableGA := siteInfo.GaAnalysisId != ""
	enableBaidu := siteInfo.BaiduAnalysisId != ""
	topViewer := a.ArticleSvr.GetTopViewer("list", viewerDataNum)
	topVisited := a.ArticleSvr.GetTopVisited("list", viewerDataNum)
	recentVisitedArticles := a.ArticleSvr.GetRecentVisitedArticles("list", viewerDataNum)

	lastVisitItem := a.VisitSvr.GetLastVisitItem()

	var LastVisitedTime time.Time
	var pathanme string
	if lastVisitItem != nil {
		LastVisitedTime = lastVisitItem.LastVisitedTime
		pathanme = lastVisitItem.Pathname
	}

	viewer := a.MetaSvr.GetViewer()
	maxArticleViewer := int64(0)
	maxArticleVisited := int64(0)

	if topViewer != nil && len(topViewer) > 0 {
		maxArticleViewer = (topViewer)[0].Viewer
	}
	if topVisited != nil && len(topVisited) > 0 {
		maxArticleVisited = (topVisited)[0].Visited
	}
	return &domain.ViewerTabData{
		EnableGA:                enableGA,
		EnableBaidu:             enableBaidu,
		TopViewer:               topViewer,
		TopVisited:              topVisited,
		RecentVisitArticles:     recentVisitedArticles,
		SiteLastVisitedTime:     LastVisitedTime,
		SiteLastVisitedPathname: pathanme,
		TotalViewer:             viewer.Viewer,
		TotalVisited:            viewer.Visited,
		MaxArticleViewer:        maxArticleViewer,
		MaxArticleVisited:       maxArticleVisited,
	}
}

func (a *AnalysisService) GetArticleTabData(articleTabDataNum int64) *domain.ArticleTabData {
	articleNum := a.ArticleSvr.GetTotalNum(true)
	wordNum := a.MetaSvr.GetTotalWords()
	tagNum := len(a.TagSvr.GetAllTags(true))
	categoryNum := len(a.CategorySvr.GetAllCategories())
	categoryPieData := a.CategorySvr.GetPipeData()
	columnData := a.TagSvr.GetColumnData(articleTabDataNum, true)
	return &domain.ArticleTabData{
		ArticleNum:      articleNum,
		WordNum:         wordNum,
		TagNum:          int64(tagNum),
		CategoryNum:     int64(categoryNum),
		CategoryPieData: categoryPieData,
		ColumnData:      columnData,
	}
}
