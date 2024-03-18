package router

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/event"
	"cc.allio/fusion/internal/svr"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

const TagPathPrefix = "/api/admin/tag"

type TagRouter struct {
	Cfg        *config.Config
	TagService *svr.TagService
	Isr        *event.IsrEventBus
}

var TagRouterSet = wire.NewSet(wire.Struct(new(TagRouter), "*"))

// GetAllTags
// @Summary get all tags
// @Schemes
// @Description get all tags
// @Tags Tag
// @Accept json
// @Produce json
// @Success 200 {object} []string
// @Router /api/admin/tag/all [Get]
func (t *TagRouter) GetAllTags(c *gin.Context) *R {
	tags := t.TagService.GetAllTags(true)
	return Ok(tags)
}

// GetArticlesByTagName
// @Summary get articles by tag name
// @Schemes
// @Description get articles by tag name
// @Tags Tag
// @Accept json
// @Produce json
// @Success 200 {object} []domain.Article
// @Router /api/admin/tag/:name [Get]
func (t *TagRouter) GetArticlesByTagName(c *gin.Context) *R {
	if t.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	tagName := c.Param("name")
	articles := t.TagService.GetArticlesByTagName(tagName, true)
	return Ok(articles)
}

// UpdateTagByName
// @Summary update articles tags by tag name
// @Schemes
// @Description update articles by tag name
// @Tags Tag
// @Accept json
// @Produce json
// @Param        value   path      string  true  "value"
// @Success 200 {object} bool
// @Router /api/admin/tag/:name [Put]
func (t *TagRouter) UpdateTagByName(c *gin.Context) *R {
	if t.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	old := c.Param("name")
	set := c.Query("value")
	result := t.TagService.UpdateArticleTag(old, set)
	t.Isr.ActiveAll("trigger incremental rendering by tags update")
	return Ok(result)
}

// DeleteTagByName
// @Summary delete articles tags by tag name
// @Schemes
// @Description delete articles tags by tag name
// @Tags Tag
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/tag/:name [Delete]
func (t *TagRouter) DeleteTagByName(c *gin.Context) *R {
	if t.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	name := c.Param("name")
	result := t.TagService.DeleteArticleTag(name)
	t.Isr.ActiveAll("trigger incremental rendering by tags delete")
	return Ok(result)
}

func (t *TagRouter) Register(r *gin.Engine) {
	r.GET(TagPathPrefix+"/all", Handle(t.GetAllTags))

	r.GET(TagPathPrefix+"/:name", Handle(t.GetArticlesByTagName))
	r.PUT(TagPathPrefix+"/:name", Handle(t.UpdateTagByName))
	r.DELETE(TagPathPrefix+"/:name", Handle(t.DeleteTagByName))
}
