package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/pkg/mongodb"
	"cc.allio/fusion/pkg/util"
	"errors"
	"github.com/google/wire"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slog"
	"strconv"
	"strings"
	"time"
)

type ArticleView = string

const (
	PublicArticleView = "public"
	AdminArticleView  = "admin"
	ListArticleView   = "list"
)

var (
	PublicView = bson.D{
		{"title", 1},
		{"content", 1},
		{"tags", 1},
		{"category", 1},
		{"updatedAt", 1},
		{"createdAt", 1},
		{"lastVisitedTime", 1},
		{"id", 1},
		{"top", 1},
		{"_id", 0},
		{"viewer", 1},
		{"visited", 1},
		{"private", 1},
		{"hidden", 1},
		{"author", 1},
		{"copyright", 1},
		{"pathname", 1},
	}

	AdminView = bson.D{
		{"title", 1},
		{"content", 1},
		{"tags", 1},
		{"category", 1},
		{"updatedAt", 1},
		{"createdAt", 1},
		{"lastVisitedTime", 1},
		{"id", 1},
		{"top", 1},
		{"_id", 0},
		{"viewer", 1},
		{"visited", 1},
		{"private", 1},
		{"hidden", 1},
		{"author", 1},
		{"copyright", 1},
		{"pathname", 1},
	}

	ListView = bson.D{
		{"title", 1},
		{"content", 1},
		{"tags", 1},
		{"category", 1},
		{"updatedAt", 1},
		{"createdAt", 1},
		{"lastVisitedTime", 1},
		{"id", 1},
		{"top", 1},
		{"_id", 0},
		{"viewer", 1},
		{"visited", 1},
		{"private", 1},
		{"hidden", 1},
		{"author", 1},
		{"copyright", 1},
		{"pathname", 1},
	}
)

type ArticleService struct {
	Cfg          *config.Config
	MetaService  *MetaService
	ArticleRepo  *repo.ArticleRepository
	CategoryRepo *repo.CategoryRepository
}

var ArticleServiceSet = wire.NewSet(wire.Struct(new(ArticleService), "*"))

func (a *ArticleService) GetByLink(link string) []*domain.Article {
	filter := mongodb.NewLogicalDefaultLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	filter.Append(bson.E{Key: "Content", Value: bson.D{{"$regex", link}, {"&options", "i"}}})
	articles, err := a.ArticleRepo.FindList(filter)
	if err != nil {
		slog.Error("Get article by link has err", "err", err)
		return make([]*domain.Article, 0)
	}
	return articles
}

func (a *ArticleService) DeleteById(id int64) bool {
	filter := mongodb.NewLogicalDefault(bson.E{Key: "id", Value: id})
	update := bson.D{{"$set", bson.D{{"deleted", true}}}}
	updated, err := a.ArticleRepo.Update(filter, update)
	if err != nil {
		slog.Error("Delete article has err", "err", err)
		return false
	}
	return updated
}

func (a *ArticleService) Create(article *domain.Article) (*domain.Article, error) {
	id, err := a.ArticleRepo.Save(article)
	if err != nil {
		return nil, err
	}
	article.Id = id
	return article, nil
}

func (a *ArticleService) UpdateById(id int64, article *domain.Article) bool {
	filter := mongodb.NewLogicalDefault(bson.E{Key: "id", Value: id})
	articleBson := util.ToBsonElements(article)
	update := bson.D{{"$set", articleBson}}
	updated, err := a.ArticleRepo.Update(filter, update)
	if err != nil {
		return false
	}
	return updated
}

func (a *ArticleService) GetById(id int64) *domain.Article {
	filter := mongodb.NewLogicalDefault(bson.E{Key: "id", Value: id})
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	article, err := a.ArticleRepo.FindOne(filter)
	if err != nil {
		return &domain.Article{}
	}
	return article
}

func (a *ArticleService) GetByPathname(pathname string) *domain.Article {
	filter := mongodb.NewLogicalDefault(bson.E{Key: "pathname", Value: pathname})
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	article, err := a.ArticleRepo.FindOne(filter)
	if err != nil {
		return &domain.Article{}
	}
	return article
}

func (a *ArticleService) GetByOption(option credential.ArticleSearchOptionCredential, isPublic bool) *domain.ArticlePageResult {
	filter := mongodb.NewLogical()
	// delete
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	if isPublic {
		// hidden
		filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(HiddenFilter))
	}
	opt := options.FindOptions{}

	// sort
	var sort bson.E
	if option.SortTop == domain.AscSort {
		sort = bson.E{Key: "top", Value: 1}
	} else if option.SortTop == domain.DescSort {
		sort = bson.E{Key: "top", Value: -1}
	}
	if option.SortViewer == domain.AscSort {
		sort = bson.E{Key: "viewer", Value: 1}
	} else if option.SortViewer == domain.DescSort {
		sort = bson.E{Key: "viewer", Value: -1}
	}
	if option.SortCreatedAt == domain.AscSort {
		sort = bson.E{Key: "created", Value: 1}
	} else if option.SortCreatedAt == domain.DescSort {
		sort = bson.E{Key: "created", Value: -1}
	}
	if &sort != nil {
		opt.Sort = sort
	}

	// tags
	if lo.IsNotEmpty(option.Tags) {
		tags := strings.Split(option.Tags, ",")
		tagFilters := make([]bson.E, 0)
		for _, tag := range tags {
			if option.RegMatch {
				tagFilters = append(tagFilters, bson.E{Key: "tags", Value: bson.D{{"$regex", tag}, {"$options", "i"}}})
			} else {
				tagFilters = append(tagFilters, bson.E{Key: "tags", Value: tag})
			}
		}
		if len(tagFilters) > 0 {
			filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(tagFilters))
		}
	}

	// category
	if lo.IsNotEmpty(option.Category) {
		if option.RegMatch {
			filter.Append(bson.E{Key: "category", Value: bson.D{{"$regex", option.Category}, {"$options", "i"}}})
		} else {
			filter.Append(bson.E{Key: "category", Value: option.Category})
		}
	}

	// title
	if lo.IsNotEmpty(option.Title) {
		filter.Append(bson.E{Key: "title", Value: bson.D{{"$regex", option.Title}, {"$options", "i"}}})
	}

	// time
	if lo.IsNotEmpty(option.StartTime) || lo.IsNotEmpty(option.EndTime) {
		timeFilter := make([]bson.E, 0)
		if lo.IsNotEmpty(option.StartTime) {
			startTime, err := time.Parse(time.DateTime, option.StartTime)
			if err == nil {
				timeFilter = append(timeFilter, bson.E{Key: "$gte", Value: startTime})
			}
		}
		if lo.IsNotEmpty(option.EndTime) {
			endTime, err := time.Parse(time.DateTime, option.EndTime)
			if err == nil {
				timeFilter = append(timeFilter, bson.E{Key: "$lte", Value: endTime})
			}
		}
		if len(timeFilter) > 0 {
			filter.Append(bson.E{Key: "createdAt", Value: timeFilter})
		}
	}

	if option.PageSize != -1 && !isPublic {
		skip := int64(option.PageSize*option.Page - option.PageSize)
		limit := int64(option.PageSize)
		opt.Skip = &skip
		opt.Limit = &limit
	}

	articles, err := a.ArticleRepo.FindList(filter, &opt)

	if err != nil {
		slog.Error("GetByOption has err", "err", err.Error())
		return &domain.ArticlePageResult{Articles: make([]*domain.Article, 0), Total: 0, TotalWordCount: 0}
	}

	// 逻辑分页
	if isPublic && option.PageSize != -1 {
		topArticles := lo.Filter(articles, func(article *domain.Article, index int) bool {
			top := article.Top
			return util.ToIntBool(int(top))
		})
		exclusionArticles := lo.Filter(articles, func(article *domain.Article, index int) bool {
			top := article.Top
			return !util.ToIntBool(int(top))
		})
		articles = util.Combination[*domain.Article](topArticles, exclusionArticles)
		skip := (option.Page - 1) & option.PageSize
		end := skip + option.PageSize
		if end > len(articles)-1 {
			end = len(articles)
		}
		articles = articles[skip:end]
	}

	count, err := a.ArticleRepo.Count(filter)

	if err != nil {
		slog.Error("GetByOption query count has err", "err", err)
		count = 0
	}

	// 过滤私有文章
	if isPublic {
		articles = lo.Map(articles, func(article *domain.Article, index int) *domain.Article {
			isPrivate := article.Private
			category, err := a.CategoryRepo.FindOne(filter.Append(bson.E{Key: "name", Value: article.Category}))
			if err != nil {
				isPrivate = false
			} else {
				isPrivate = isPrivate && category.Private
			}
			if isPrivate {
				return article
			} else {
				return &domain.Article{
					Id:              article.Id,
					Title:           article.Title,
					Content:         "",
					Tags:            article.Tags,
					Top:             article.Top,
					Category:        article.Category,
					Hidden:          article.Hidden,
					Author:          article.Author,
					Pathname:        article.Pathname,
					Private:         true,
					Password:        "",
					Deleted:         article.Deleted,
					Viewer:          article.Viewer,
					Visited:         article.Visited,
					Copyright:       article.Copyright,
					LastVisitedTime: article.LastVisitedTime,
					CreatedAt:       article.CreatedAt,
					UpdatedAt:       article.UpdatedAt,
				}
			}
		})
	}

	// 返回结果
	result := &domain.ArticlePageResult{}
	totalWordCount := int64(0)
	result.TotalWordCount = totalWordCount
	if option.WithWordCount {
		totalWordCount = lo.Reduce(
			articles,
			func(total int64, article *domain.Article, index int) int64 {
				total += util.WordCount(article.Content)
				return total
			},
			int64(0))
	}
	if option.WithWordCount && option.ToListView {
		// 重置视图
		articles = lo.Map(articles, func(article *domain.Article, index int) *domain.Article {
			return &domain.Article{
				Id:              article.Id,
				Title:           article.Title,
				Content:         "",
				Tags:            article.Tags,
				Top:             article.Top,
				Category:        article.Category,
				Hidden:          article.Hidden,
				Author:          article.Author,
				Pathname:        article.Pathname,
				Private:         article.Private,
				Password:        "",
				Deleted:         article.Deleted,
				Viewer:          article.Viewer,
				Visited:         article.Visited,
				Copyright:       article.Copyright,
				LastVisitedTime: article.LastVisitedTime,
				CreatedAt:       article.CreatedAt,
				UpdatedAt:       article.UpdatedAt,
			}
		})
	}
	result.Articles = articles
	result.Total = count
	return result
}

func (a *ArticleService) GetAll(view ArticleView, includeHidden bool, includeDelete bool) []*domain.Article {
	viewFilter := a.getView(view)
	if viewFilter == nil {
		slog.Error("view incorrect")
		return []*domain.Article{}
	}
	filter := mongodb.NewLogicalDefaultArray(viewFilter)
	orLogic := mongodb.NewLogicalOr()
	if !includeDelete {
		orLogic = orLogic.AppendArray(DeleteFilter)
	}
	if !includeHidden {
		orLogic = orLogic.AppendArray(HiddenFilter)
	}
	filter.AppendLogical(orLogic)

	sort := mongodb.MakeSort("createdAt", mongodb.Descending)
	opt := options.FindOptions{Sort: sort}
	articles, err := a.ArticleRepo.FindList(filter, &opt)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}
	return articles
}

func (a *ArticleService) GetRecentVisitedArticles(view ArticleView, limit int64) []*domain.Article {
	viewFilter := a.getView(view)
	if viewFilter == nil {
		slog.Error("view incorrect")
		return []*domain.Article{}
	}
	filter := mongodb.NewLogicalDefaultArray(viewFilter)
	filter.Append(bson.E{Key: "lastVisitedTime", Value: bson.E{Key: "$exists", Value: true}})
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	opt := options.FindOptions{Sort: bson.E{Key: "lastVisitedTime", Value: -1}, Limit: &limit}
	articles, err := a.ArticleRepo.FindList(filter, &opt)
	if err != nil {
		slog.Error(err.Error())
		return []*domain.Article{}
	}
	return articles
}

func (a *ArticleService) GetTopViewer(view ArticleView, limit int64) []*domain.Article {
	viewFilter := a.getView(view)
	if viewFilter == nil {
		slog.Error("view incorrect")
		return []*domain.Article{}
	}
	filter := mongodb.NewLogicalDefaultArray(viewFilter)
	filter.Append(bson.E{Key: "viewer", Value: bson.D{{"$ne", 0}, {"$exists", true}}})
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	opt := options.FindOptions{Sort: bson.E{Key: "viewer", Value: -1}, Limit: &limit}
	articles, err := a.ArticleRepo.FindList(filter, &opt)
	if err != nil {
		slog.Error(err.Error())
		return []*domain.Article{}
	}
	return articles
}

func (a *ArticleService) GetTopVisited(view ArticleView, limit int64) []*domain.Article {
	viewFilter := a.getView(view)
	if viewFilter == nil {
		slog.Error("view incorrect")
		return []*domain.Article{}
	}
	filter := mongodb.NewLogicalDefaultArray(viewFilter)
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	filter.Append(bson.E{Key: "viewer", Value: bson.D{{"$ne", 0}, {"$exists", true}}})
	opt := options.FindOptions{Sort: bson.E{Key: "visited", Value: -1}, Limit: &limit}
	articles, err := a.ArticleRepo.FindList(filter, &opt)
	if err != nil {
		slog.Error(err.Error())
		return []*domain.Article{}
	}
	return articles
}

func (a *ArticleService) getView(view ArticleView) bson.D {
	switch view {
	case AdminArticleView:
		return AdminView
	case PublicArticleView:
		return PublicView
	case ListArticleView:
		return ListView
	}
	return nil
}

func (a *ArticleService) GetTotalNum(includeHidden bool) int64 {
	or := mongodb.NewLogicalOr()
	or = or.AppendArray(DeleteFilter)
	if includeHidden {
		or = or.AppendArray(HiddenFilter)
	}
	count, err := a.ArticleRepo.Count(or)
	if err != nil {
		slog.Error("GetTotalNum has err", "err", err.Error())
		return 0
	}
	return count
}

// UpdateTags by article id set new tags
func (a *ArticleService) UpdateTags(id int64, newTags []string) (bool, error) {
	return a.ArticleRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: id}), bson.D{{"tags", newTags}})
}

func (a *ArticleService) getArticleByIdOrPathname(idOrPathname string) *domain.Article {
	var article *domain.Article
	id, err := strconv.Atoi(idOrPathname)
	if err == nil {
		article = a.GetById(int64(id))
	} else {
		article = a.GetByPathname(idOrPathname)
	}
	return article
}

func (a *ArticleService) GetArticleByIdOrPathnameWithAlternate(idOrPathname string) (*credential.AlternateArticle, error) {
	article := a.getArticleByIdOrPathname(idOrPathname)
	if article == nil {
		return nil, errors.New("not found article")
	}

	if article.Hidden {
		siteInfo := a.MetaService.GetSiteInfo()
		if siteInfo != nil && siteInfo.AllowOpenHiddenPostByUrl == "false" {
			return nil, errors.New("article invisible")
		}
	}

	if article.Private {
		article.Content = ""
	} else {
		category, err := a.CategoryRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "name", Value: article.Category}))
		if err == nil && category != nil && category.Private {
			article.Private = true
			article.Content = ""
		}
	}
	preArticle := a.getPreArticle(article.CreatedAt, false)
	nextArticle := a.getNextArticle(article.CreatedAt, false)
	return &credential.AlternateArticle{Article: article, Pre: preArticle, Next: nextArticle}, nil
}

func (a *ArticleService) getPreArticle(createTime time.Time, includeHidden bool) *domain.Article {
	filter := mongodb.NewLogical()
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	filter.Append(bson.E{Key: "createdAt", Value: bson.E{Key: "$lt", Value: createTime}})
	if includeHidden {
		filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(HiddenFilter))
	}
	limit := int64(-1)
	option := &options.FindOptions{Sort: bson.E{Key: "createAt", Value: -1}, Limit: &limit}
	articles, err := a.ArticleRepo.FindList(filter, option)
	if err != nil {
		slog.Error("find pre article has error", "err", err, "time", createTime)
		return nil
	}
	if len(articles) > 0 {
		return articles[0]
	} else {
		return nil
	}
}

func (a *ArticleService) getNextArticle(createTime time.Time, includeHidden bool) *domain.Article {
	filter := mongodb.NewLogical()
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	filter.Append(bson.E{Key: "createdAt", Value: bson.E{Key: "$gt", Value: createTime}})
	if includeHidden {
		filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(HiddenFilter))
	}
	limit := int64(-1)
	option := &options.FindOptions{Sort: bson.E{Key: "createAt", Value: -1}, Limit: &limit}
	articles, err := a.ArticleRepo.FindList(filter, option)
	if err != nil {
		slog.Error("find pre article has error", "err", err, "time", createTime)
		return nil
	}
	if len(articles) > 0 {
		return articles[0]
	} else {
		return nil
	}
}

func (a *ArticleService) GetArticleByIdOrPathnameWithPassword(idOrPathname string, password string) *domain.Article {
	if lo.IsEmpty(password) {
		return nil
	}
	article := a.getArticleByIdOrPathname(idOrPathname)
	if article == nil {
		return nil
	}
	category, err := a.CategoryRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "name", Value: article.CreatedAt}))

	var categoryPassword string
	if err == nil && category != nil && category.Private {
		categoryPassword = category.Password
	}
	var targetPassword string
	if lo.IsNotEmpty(categoryPassword) {
		targetPassword = categoryPassword
	} else {
		targetPassword = article.Password
	}
	if lo.IsEmpty(targetPassword) {
		article.Password = ""
		return article
	}
	if targetPassword == password {
		article.Password = ""
		return article
	}
	return nil
}

func (a *ArticleService) SearchByText(text string, includeHidden bool) []*domain.Article {
	filter := mongodb.NewLogical()
	regexText := bson.D{{"$regex", text}, {"$options", "i"}}
	filter.AppendLogical(mongodb.NewLogicalDefaultArray(bson.D{{"content", regexText}, {"title", regexText}, {"category", regexText}, {"tags", regexText}}))
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	if includeHidden {
		filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(HiddenFilter))
	}
	articles, err := a.ArticleRepo.FindList(filter)
	if err != nil {
		return make([]*domain.Article, 0)
	}
	return articles
}

func (a *ArticleService) UpdateViewerByPathname(pathname string, isNew bool) {
	article := a.getArticleByIdOrPathname(pathname)
	if article != nil {
		oldViewer := article.Visited
		oldVisited := article.Visited
		newViewer := oldViewer + 1
		newVisited := oldVisited
		if isNew {
			newVisited = oldVisited + 1
		}
		a.ArticleRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: article.Id}), bson.D{{"visited", newVisited}, {"viewer", newViewer}})
	}
}

func (a *ArticleService) GetTimeLine() map[string][]*domain.Article {
	filter := mongodb.NewLogical()
	// append default domain model
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(HiddenFilter))
	// sort
	sort := mongodb.MakeSort("createdAt", mongodb.Descending)
	opt := &options.FindOptions{Sort: sort}
	articles, err := a.ArticleRepo.FindList(filter, opt)
	if err != nil {
		return map[string][]*domain.Article{}
	}
	return lo.GroupBy[*domain.Article, string](articles, func(article *domain.Article) string {
		return string(rune(article.CreatedAt.Year()))
	})
}
