package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/web"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

const AnalysisPathPrefix = "/api/admin/analysis"

type AnalysisRoute struct {
	Cfg         *config.Config
	AnalysisSvr *svr.AnalysisService
}

var AnalysisRouterSet = wire.NewSet(wire.Struct(new(AnalysisRoute), "*"))

// GetWelcomePageData
// @Summary 获取首页统计数据
// @Schemes
// @Description statistic home page data
// @Tags Analysis
// @Accept json
// @Produce json
// @Param        tab                query      string      true       "tab"              Enums(overview, viewer, article)
// @Param        viewerDataNum      query      int         false      "overviewDataNum"
// @Param        overviewDataNum    query      int         false      "overviewDataNum"
// @Param        articleTabDataNum  query      int         false      "articleTabDataNum"
// @Router /api/admin/analysis [Get]
func (a *AnalysisRoute) GetWelcomePageData(c *gin.Context) *R {
	tab := c.Query("tab")
	if tab == "overview" {
		viewerDataNum := web.ParseNumberForQuery(c, "viewerDataNum", 0)
		overviewData := a.AnalysisSvr.GetOverViewTabData(int64(viewerDataNum))
		return Ok(overviewData)
	} else if tab == "viewer" {
		overviewDataNum := web.ParseNumberForQuery(c, "overviewDataNum", 0)
		viewerTabData := a.AnalysisSvr.GetViewerTabData(int64(overviewDataNum))
		return Ok(viewerTabData)
	} else if tab == "article" {
		articleTabDataNum := web.ParseNumberForQuery(c, "articleTabDataNum", 0)
		articleTabData := a.AnalysisSvr.GetArticleTabData(int64(articleTabDataNum))
		return Ok(articleTabData)
	} else {
		return InternalError(errors.New("not found any tab view data"))
	}
}

func (a *AnalysisRoute) Register(r *gin.Engine) {
	r.GET(AnalysisPathPrefix, Handle(a.GetWelcomePageData))
}
