package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/svr"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

const LinkPathPrefix = "/api/admin/meta/link"

type LinkRoute struct {
	Cfg         *config.Config
	MetaService *svr.MetaService
}

var LinkRouterSet = wire.NewSet(wire.Struct(new(LinkRoute), "*"))

// GetLink
// @Summary get link by meta
// @Schemes
// @Description get link by meta
// @Tags Link
// @Accept json
// @Produce json
// @Success 200 {object} []domain.LinkItem
// @Router /api/admin/meta/link [Get]
func (l *LinkRoute) GetLink(c *gin.Context) *R {
	links := l.MetaService.GetLinks()
	return Ok(links)
}

// UpdateLink
// @Summary update link
// @Schemes
// @Description update link
// @Tags Link
// @Accept json
// @Produce json
// @Param        link   body      domain.LinkItem   true  "link"
// @Success 200 {object} bool
// @Router /api/admin/meta/link [Put]
func (l *LinkRoute) UpdateLink(c *gin.Context) *R {
	if l.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	link := &domain.LinkItem{}
	if err := c.Bind(link); err != nil {
		return InternalError(err)
	}
	updated, err := l.MetaService.AddOrUpdateLink(link)
	if err != nil {
		return InternalError(err)
	}
	return Ok(updated)
}

// CreateLink
// @Summary create link
// @Schemes
// @Description create link
// @Tags Link
// @Accept json
// @Produce json
// @Param        link   body      domain.LinkItem   true  "link"
// @Success 200 {object} bool
// @Router /api/admin/meta/link [Post]
func (l *LinkRoute) CreateLink(c *gin.Context) *R {
	if l.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	link := &domain.LinkItem{}
	if err := c.Bind(link); err != nil {
		return InternalError(err)
	}
	added, err := l.MetaService.AddOrUpdateLink(link)
	if err != nil {
		return InternalError(err)
	}
	return Ok(added)
}

// DeleteLink
// @Summary delete link by link name
// @Schemes
// @Description delete link by link name
// @Tags Link
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/meta/link/:name [Delete]
func (l *LinkRoute) DeleteLink(c *gin.Context) *R {
	if l.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	name := c.Query("name")
	deleted, err := l.MetaService.DeleteLinkByName(name)
	if err != nil {
		return InternalError(err)
	}
	return Ok(deleted)
}

func (l *LinkRoute) Register(r *gin.Engine) {
	r.GET(LinkPathPrefix, Handle(l.GetLink))
	r.PUT(LinkPathPrefix, Handle(l.UpdateLink))
	r.POST(LinkPathPrefix, Handle(l.CreateLink))
	r.DELETE(LinkPathPrefix+"/:name", Handle(l.DeleteLink))
}
