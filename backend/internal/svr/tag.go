package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/domain"
	"github.com/google/wire"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
)

type TagService struct {
	Cfg        *config.Config
	ArticleSvr *ArticleService
}

var TagServiceSet = wire.NewSet(wire.Struct(new(TagService), "*"))

func (t *TagService) GetColumnData(topNum int64, includeHidden bool) []*domain.TypeValue[int64] {
	data := t.GetTagsWithArticle(includeHidden)
	tags := make([]string, len(data))
	for key := range data {
		tags = append(tags, key)
	}
	typeValues := make([]*domain.TypeValue[int64], 0)
	for index, tag := range tags {
		if index+1 == int(topNum) {
			break
		}
		typeValues = append(typeValues, &domain.TypeValue[int64]{
			Type:  tag,
			Value: int64(len(data[tag])),
		})
	}
	return typeValues
}

func (t *TagService) GetAllTags(includeHidden bool) []string {
	allArticles := t.GetTagsWithArticle(includeHidden)
	tags := make([]string, 0)
	for k, _ := range allArticles {
		tags = append(tags, k)
	}
	return tags
}

func (t *TagService) GetArticlesByTagName(tagName string, includeHidden bool) []*domain.Article {
	allArticles := t.GetTagsWithArticle(includeHidden)
	return allArticles[tagName]
}

func (t *TagService) GetTagsWithArticle(includeHidden bool) map[string][]*domain.Article {
	allArticles := t.ArticleSvr.GetAll("list", includeHidden, false)
	data := make(map[string][]*domain.Article)
	for _, article := range allArticles {
		for _, tag := range article.Tags {
			if _, ok := data[tag]; !ok {
				data[tag] = make([]*domain.Article, 0)
			}
			data[tag] = append(data[tag], article)
		}
	}
	return data
}

func (t *TagService) UpdateArticleTag(old string, new string) *credential.TagResultCredential {
	articles := t.GetArticlesByTagName(old, true)
	for _, article := range articles {
		// filter old tag
		tags := lo.Filter(article.Tags, func(tag string, index int) bool { return tag != old })
		// add new tag
		tags = append(tags, new)
		// update article, ignore error
		_, err := t.ArticleSvr.UpdateTags(article.Id, tags)
		if err != nil {
			slog.Error("update article tags has error", "err", err, "tags", tags)
		}
	}
	return &credential.TagResultCredential{Message: "更新成功", Total: len(articles)}
}

func (t *TagService) DeleteArticleTag(tagName string) *credential.TagResultCredential {
	articles := t.GetArticlesByTagName(tagName, true)
	for _, article := range articles {
		// filter old tag
		tags := lo.Filter(article.Tags, func(tag string, index int) bool { return tag != tagName })
		// update article, ignore error
		_, err := t.ArticleSvr.UpdateTags(article.Id, tags)
		if err != nil {
			slog.Error("update article tags has error", "err", err, "tags", tags)
		}
	}
	return &credential.TagResultCredential{Message: "删除成功", Total: len(articles)}
}
