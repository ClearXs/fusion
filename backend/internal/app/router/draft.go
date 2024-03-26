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
	"net/http"
)

const DraftPathPrefix = "/api/admin/draft"

type DraftRoute struct {
	Cfg          *config.Config
	DraftService *svr.DraftService
	Isr          *event.IsrEventBus
	Script       *event.ScriptEngine
}

var DraftRouterSet = wire.NewSet(wire.Struct(new(DraftRoute), "*"))

// GetDraftByOption
// @Summary get drafts by options
// @Schemes
// @Description get drafts by options
// @Tags Draft
// @Accept json
// @Produce json
// @Param        page               query      int        false       "page"               -1
// @Param        pageSize           query      int        false       "pageSize"            5
// @Param        toListView         query      bool       false       "toListView"          false
// @Param        category           query      string     false       "category"
// @Param        tags               query      string     false       "page"
// @Param        title              query      int        false       "title"
// @Param        sortCreateAt       query      int        false       "sortCreateAt"
// @Param        startTime          query      int        false       "startTime"
// @Param        endTime            query      int        false       "endTime"
// @Success 200 {object} domain.DraftPageResult
// @Router /api/admin/draft [Get]
func (d *DraftRoute) GetDraftByOption(c *gin.Context) *R {
	page := web.ParseNumberForQuery(c, "page", -1)
	pageSize := web.ParseNumberForQuery(c, "pageSize", 5)
	toListView := web.ParseBoolForQuery(c, "toListView", false)
	category := c.Query("category")
	tags := c.Query("tags")
	title := c.Query("title")
	sortCreateAt := c.Query("sortCreateAt")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	option := &credential.DraftSearchOptionCredential{
		Page:          page,
		PageSize:      pageSize,
		ToListView:    toListView,
		Category:      category,
		Tags:          tags,
		Title:         title,
		SortCreatedAt: sortCreateAt,
		StartTime:     startTime,
		EndTime:       endTime,
	}
	result := d.DraftService.GetByOption(option)
	return Ok(result)
}

// GetDraft
// @Summary get draft by id
// @Schemes
// @Description get draft by id
// @Tags Draft
// @Accept json
// @Produce json
// @Success 200 {object} domain.Draft
// @Router /api/admin/draft/:id [Get]
func (d *DraftRoute) GetDraft(c *gin.Context) *R {
	id := web.ParseNumberForQuery(c, "id", -1)
	draft, err := d.DraftService.GetById(int64(id))
	if err != nil {
		return InternalError(err)
	}
	return Ok(draft)
}

// UpdateDraft
// @Summary update draft by id
// @Schemes
// @Description update draft by id
// @Tags Draft
// @Accept json
// @Produce json
// @Param        draft   body      domain.Draft   true  "draft"
// @Success 200 {object} bool
// @Router /api/admin/draft/:id [Put]
func (d *DraftRoute) UpdateDraft(c *gin.Context) *R {
	draft := &domain.Draft{}
	if err := c.Bind(draft); err != nil {
		return InternalError(err)
	}
	id := web.ParseNumberForQuery(c, "id", -1)
	d.Script.DispatchBeforeUpdateDraftEvent(draft, id)
	successed, err := d.DraftService.UpdateById(int64(id), draft)
	if err != nil {
		return InternalError(err)
	}
	d.Script.DispatchAfterUpdateDraftEvent(draft, id, successed)
	return Ok(successed)
}

// CreateDraft
// @Summary crate draft
// @Schemes
// @Description crate draft
// @Tags Draft
// @Accept json
// @Produce json
// @Param        draft   body      domain.Draft   true  "draft"
// @Success 200 {object} domain.Draft
// @Router /api/admin/draft [Post]
func (d *DraftRoute) CreateDraft(c *gin.Context) *R {
	draft := &domain.Draft{}
	if err := c.Bind(draft); err != nil {
		return InternalError(err)
	}
	d.Script.DispatchBeforeUpdateDraftEvent(draft)
	newDraft, err := d.DraftService.Create(draft)
	if err != nil {
		return InternalError(err)
	}
	d.Script.DispatchAfterUpdateDraftEvent(draft)
	return Ok(newDraft)
}

// DeleteDraft
// @Summary delete draft by id
// @Schemes
// @Description crate draft
// @Tags Draft
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/draft/:id [Delete]
func (d *DraftRoute) DeleteDraft(c *gin.Context) *R {
	id := web.ParseNumberForQuery(c, "id", -1)
	successed, err := d.DraftService.DeleteById(int64(id))
	if err != nil {
		return InternalError(err)
	}
	d.Script.DispatchDeleteDraftEvent(id)
	return Ok(successed)
}

// PublishDraft
// @Summary publish and update draft
// @Schemes
// @Description publish and update draft
// @Tags Draft
// @Accept json
// @Produce json
// @Param        draft   body      domain.Draft   true  "draft"
// @Success 200 {object} domain.Draft
// @Router /api/admin/draft/publish/:id [Post]
func (d *DraftRoute) PublishDraft(c *gin.Context) *R {
	if d.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止发布草稿！"))
	}
	id := web.ParseNumberForQuery(c, "id", -1)
	option := &credential.DraftPublishCredential{}
	if err := c.Bind(option); err != nil {
		return InternalError(err)
	}
	d.Script.DispatchBeforeUpdateDraftEvent(option)
	newDraft, err := d.DraftService.Publish(int64(id), option)
	if err != nil {
		return InternalError(err)
	}
	d.Isr.ActiveAll("trigger incremental rendering by public draft")
	d.Script.DispatchAfterUpdateDraftEvent(newDraft)
	return Ok(newDraft)
}

func (d *DraftRoute) Register(r *gin.Engine) {
	r.GET(DraftPathPrefix, Handle(d.GetDraftByOption))

	r.GET(DraftPathPrefix+"/:id", Handle(d.GetDraft))
	r.PUT(DraftPathPrefix+"/:id", Handle(d.UpdateDraft))
	r.POST(DraftPathPrefix+"/:id", Handle(d.CreateDraft))
	r.DELETE(DraftPathPrefix+"/:id", Handle(d.DeleteDraft))

	r.POST(DraftPathPrefix+"/publish/:id", Handle(d.PublishDraft))
}
