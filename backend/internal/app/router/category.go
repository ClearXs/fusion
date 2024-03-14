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
	"net/http"
)

const CategoryPathPrefix = "/api/admin/category"

type CategoryRouter struct {
	Cfg             *config.Config
	CategoryService *svr.CategoryService
}

var CategoryRouterSet = wire.NewSet(wire.Struct(new(CategoryRouter), "*"))

// GetAllTags
// @Summary get all category tags
// @Schemes
// @Description get all category tags
// @Tags Category
// @Accept json
// @Produce json
// @Success 200 {object} nil
// @Router /api/admin/category/all [Get]
func (category *CategoryRouter) GetAllTags(c *gin.Context) *R {
	withDetails := web.ParseBoolForQuery(c, "detail", false)
	if withDetails {
		categories := category.CategoryService.GetAllCategories()
		return Ok(categories)
	} else {
		keys := category.CategoryService.GetAllCategoryKeys()
		return Ok(keys)
	}
}

// GetArticlesByName
// @Summary get articles by category name
// @Schemes
// @Description get articles by category name
// @Tags Category
// @Accept json
// @Produce json
// @Success 200 {object} []domain.Article
// @Router /api/admin/category/:name [Get]
func (category *CategoryRouter) GetArticlesByName(c *gin.Context) *R {
	name := c.Param("name")
	articles := category.CategoryService.GetArticlesByCategory(name, true)
	return Ok(articles)
}

// CreateCategory
// @Summary create category
// @Schemes
// @Description create category
// @Tags Category
// @Accept json
// @Produce json
// @Param        credential   body      credential.CategoryCredential   true  "credential"
// @Success 200 {object} bool
// @Router /api/admin/category [Post]
func (category *CategoryRouter) CreateCategory(c *gin.Context) *R {
	if category.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	cc := &credential.CategoryCredential{}
	if err := c.Bind(cc); err != nil {
		return InternalError(err)
	}
	successed, err := category.CategoryService.Add(cc.Name)
	if err != nil {
		return InternalError(err)
	}
	return Ok(successed)
}

// DeleteCategory
// @Summary delete category by category name
// @Schemes
// @Description delete category by category name
// @Tags Category
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/category/:name [Delete]
func (category *CategoryRouter) DeleteCategory(c *gin.Context) *R {
	if category.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	cc := &credential.CategoryCredential{}
	if err := c.Bind(cc); err != nil {
		return InternalError(err)
	}
	removed, err := category.CategoryService.Remove(cc.Name)
	if err != nil {
		return InternalError(err)
	}
	return Ok(removed)
}

// UpdateCategory
// @Summary update category by category name
// @Schemes
// @Description delete category by category name
// @Tags Category
// @Param        credential   body      domain.Category   true  "credential"
// @Accept json
// @Produce json
// @Success 200 {object} bool
// @Router /api/admin/category/:name [Put]
func (category *CategoryRouter) UpdateCategory(c *gin.Context) *R {
	if category.Cfg.Demo {
		return Error(http.StatusUnauthorized, errors.New("演示站禁止修改此项！"))
	}
	cc := &domain.Category{}
	if err := c.Bind(cc); err != nil {
		return InternalError(err)
	}
	updated, err := category.CategoryService.Update(cc)
	if err != nil {
		return InternalError(err)
	}
	return Ok(updated)
}

func (category *CategoryRouter) Register(r *gin.Engine) {
	r.GET(CategoryPathPrefix+"all", Handle(category.GetAllTags))

	r.GET(CategoryPathPrefix+"/:name", Handle(category.GetArticlesByName))

	r.POST(CategoryPathPrefix, Handle(category.CreateCategory))
	r.DELETE(CategoryPathPrefix+"/:name", Handle(category.DeleteCategory))
	r.PUT(CategoryPathPrefix+"/:name", Handle(category.UpdateCategory))
}
