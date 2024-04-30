package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/pkg/mongodb"
	"cc.allio/fusion/pkg/util"
	"errors"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slog"
)

type CategoryService struct {
	Cfg          *config.Config
	CategoryRepo *repo.CategoryRepository
	ArticleSvr   *ArticleService
}

var CategoryServiceSet = wire.NewSet(wire.Struct(new(CategoryService), "*"))

func (c *CategoryService) GetPipeData() []*domain.TypeValue[int64] {
	oldData := c.GetCategoriesWithArticle(true)
	categoryKeys := make([]string, 0)
	for key := range oldData {
		categoryKeys = append(categoryKeys, key)
	}
	if len(categoryKeys) < 0 {
		return []*domain.TypeValue[int64]{}
	}
	typeValues := make([]*domain.TypeValue[int64], 0)
	for _, key := range categoryKeys {
		if data, ok := oldData[key]; ok {
			typeValues = append(typeValues, &domain.TypeValue[int64]{
				Type:  key,
				Value: int64(len(data)),
			})
		}
	}
	return typeValues
}

func (c *CategoryService) GetAllCategories() []*domain.Category {
	categories, err := c.CategoryRepo.FindList(mongodb.NewLogical())
	if err != nil {
		slog.Error(err.Error())
		return []*domain.Category{}
	}
	return categories
}

func (c *CategoryService) GetAllCategoryKeys() []string {
	categories, err := c.CategoryRepo.FindList(mongodb.NewLogical())
	if err != nil {
		slog.Error(err.Error())
		return []string{}
	}
	if categories != nil {
		keys := make([]string, 0)

		for _, category := range categories {
			keys = append(keys, category.Name)
		}
		return keys
	}
	return []string{}
}

func (c *CategoryService) GetCategoriesWithArticle(includeHidden bool) map[string][]*domain.Article {
	allArticles := c.ArticleSvr.GetAll("list", includeHidden, false)
	categoryKeys := c.GetAllCategoryKeys()
	data := make(map[string][]*domain.Article)
	for _, key := range categoryKeys {
		data[key] = make([]*domain.Article, 0)
	}
	for _, article := range allArticles {
		category := article.Category
		if data[category] != nil {
			data[category] = append(data[category], article)
		}
	}
	return data
}

func (c *CategoryService) GetArticlesByCategory(name string, includeHidden bool) []*domain.Article {
	withArticles := c.GetCategoriesWithArticle(includeHidden)
	if articles, ok := withArticles[name]; ok {
		return articles
	} else {
		return make([]*domain.Article, 0)
	}
}

func (c *CategoryService) Add(name string) (bool, error) {
	result, err := c.CategoryRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "name", Value: name}))
	if err != nil {
		return false, nil
	}
	if result != nil {
		return false, errors.New("duplicate classification name, create failed")
	}
	if err != nil {
		return false, err
	}
	category := domain.Category{
		Name:    name,
		Type:    domain.CategoryCategoryType,
		Private: false,
	}
	successd, err := c.CategoryRepo.Save(&category)
	if err != nil {
		return false, err
	}
	return successd > 0, nil
}

func (c *CategoryService) Remove(name string) (bool, error) {
	articles := c.GetArticlesByCategory(name, true)
	if len(articles) > 0 {
		return false, errors.New("special classification contains article, can't delete operate")
	}
	removed, err := c.CategoryRepo.Remove(mongodb.NewLogicalDefault(bson.E{Key: "name", Value: name}))
	if err != nil {
		return false, err
	}
	return removed, nil
}

func (c *CategoryService) Update(entity *domain.Category) (bool, error) {
	result, err := c.CategoryRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "name", Value: entity.Name}))
	if err != nil {
		return false, nil
	}
	if result != nil {
		return false, errors.New("duplicate classification name, create failed")
	}
	// TODO 添加事物
	// 修改文章分裂
	articles := c.GetArticlesByCategory(entity.Name, true)
	for _, article := range articles {
		article.Category = entity.Name
		c.ArticleSvr.UpdateById(article.Id, article)
	}
	// 修改分类
	categoryBson := util.ToBsonElements(entity)
	update := bson.D{{"$set", categoryBson}}
	updated, err := c.CategoryRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "name", Value: entity.Name}), update)
	if err != nil {
		return false, err
	}
	return updated, nil
}
