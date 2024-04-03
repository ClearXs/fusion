package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/pkg/mongodb"
	"cc.allio/fusion/pkg/util"
	"errors"
	"github.com/google/wire"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slog"
	"strconv"
	"strings"
	"time"
)

type MetaService struct {
	Cfg           *config.Config
	MetaRepo      *repo.MetaRepository
	ViewerService *ViewerService
	VisitRepo     *repo.VisitRepository
	ArticleRepo   *repo.ArticleRepository
}

var MetaServiceSet = wire.NewSet(wire.Struct(new(MetaService), "*"))

func (m *MetaService) GetViewer() *domain.DataViewer {
	about := m.GetMeta()
	if about != nil {
		viewer := about.Viewer
		visited := about.Visited
		return &domain.DataViewer{Viewer: viewer, Visited: visited}
	}
	return &domain.DataViewer{Viewer: 0, Visited: 0}
}

func (m *MetaService) GetSiteInfo() *domain.SiteInfo {
	return m.GetMeta().SiteInfo
}

func (m *MetaService) UpdateSiteInfo(siteInfo *domain.SiteInfo) (bool, error) {
	meta := m.GetMeta()
	if meta == nil {
		return false, errors.New("meta is empty")
	}
	value := util.Composition[domain.SiteInfo](meta.SiteInfo, siteInfo)

	elements := util.ToBsonElements(value)
	return m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "_id", Value: meta.Id}), bson.D{{"$set", bson.D{{"siteInfo", elements}}}})
}

func (m *MetaService) GetTotalWords() int64 {
	about := m.GetMeta()
	if about != nil {
		return about.TotalWordCount
	}
	return 0
}

func (m *MetaService) UpdateAboutContent(content string) bool {
	succeed, err := m.MetaRepo.Update(mongodb.NewLogical(),
		bson.D{
			{"about", bson.D{{"updateAt", time.Now()}, {"content", content}}}})
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	return succeed
}

func (m *MetaService) GetMeta() *domain.Meta {
	meta, err := m.MetaRepo.FindOne(mongodb.NewLogical())
	if err != nil {
		slog.Error(err.Error())
		return &domain.Meta{}
	}
	return meta
}

// --------------------- links ---------------------

func (m *MetaService) GetLinks() []*domain.LinkItem {
	about := m.GetMeta()
	return about.Links
}

func (m *MetaService) AddOrUpdateLink(link *domain.LinkItem) (bool, error) {
	meta := m.GetMeta()
	link.UpdatedAt = time.Now()
	links := meta.Links
	if links == nil {
		links = make([]*domain.LinkItem, 0)
	}
	if _, ok := lo.Find(meta.Links, func(item *domain.LinkItem) bool { return item.Name == link.Name }); !ok {
		links = append(links, link)
	}
	return m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "_id", Value: meta.Id}), bson.D{{"$set", bson.D{{"links", links}}}})
}

func (m *MetaService) DeleteLinkByName(name string) (bool, error) {
	meta := m.GetMeta()
	if meta.Links == nil {
		return false, nil
	}
	links := meta.Links
	links = lo.Filter(links, func(item *domain.LinkItem, index int) bool { return item.Name == name })
	return m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "_id", Value: meta.Id}), bson.D{{"$set", bson.D{{"links", links}}}})
}

// --------------------- reward ---------------------

func (m *MetaService) GetRewards() []*domain.RewardItem {
	meta := m.GetMeta()
	return meta.Rewards
}

func (m *MetaService) AddOrUpdateReward(reward *domain.RewardItem) (bool, error) {
	meta := m.GetMeta()
	if meta == nil {
		return false, errors.New("meta is empty")
	}
	reward.UpdatedAt = time.Now()
	rewards := meta.Rewards
	if rewards == nil {
		rewards = make([]*domain.RewardItem, 0)
	}
	if _, ok := lo.Find(meta.Rewards, func(item *domain.RewardItem) bool { return item.Name == reward.Name }); !ok {
		rewards = append(rewards, reward)
	}
	return m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "_id", Value: meta.Id}), bson.D{{"$set", bson.D{{"$set", bson.D{{"rewards", rewards}}}}}})
}

func (m *MetaService) DeleteRewardByName(name string) (bool, error) {
	meta := m.GetMeta()
	if meta == nil {
		return false, errors.New("meta is empty")
	}
	if meta.Rewards == nil {
		return false, nil
	}
	rewards := meta.Rewards
	rewards = lo.Filter(rewards, func(item *domain.RewardItem, index int) bool { return item.Name == name })
	return m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "_id", Value: meta.Id}), bson.D{{"$set", bson.D{{"rewards", rewards}}}})
}

// --------------------- social ---------------------

func (m *MetaService) GetSocials() []*domain.SocialItem {
	meta := m.GetMeta()
	if meta == nil {
		return make([]*domain.SocialItem, 0)
	}
	return meta.Socials
}

func (m *MetaService) GetDefaultSocials() []*domain.SocialItem {
	return domain.DefaultSocials
}

func (m *MetaService) SaveOrUpdateSocial(social *domain.SocialItem) (bool, error) {
	meta := m.GetMeta()
	if meta == nil {
		return false, errors.New("meta is empty")
	}
	social.UpdatedAt = time.Now()
	socials := meta.Socials
	if socials == nil {
		socials = make([]*domain.SocialItem, 0)
	}
	if _, ok := lo.Find(meta.Socials, func(item *domain.SocialItem) bool { return item.Type == social.Type }); !ok {
		socials = append(socials, social)
	}
	return m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "_id", Value: meta.Id}), bson.D{{"$set", bson.D{{"socials", socials}}}})
}

func (m *MetaService) DeleteSocialByType(typeName string) (bool, error) {
	meta := m.GetMeta()
	if meta == nil {
		return false, errors.New("meta is empty")
	}
	if meta.Socials == nil {
		return false, nil
	}
	socials := meta.Socials
	socials = lo.Filter(socials, func(item *domain.SocialItem, index int) bool { return item.Type == typeName })
	return m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "_id", Value: meta.Id}), bson.D{{"$set", bson.D{{"socials", socials}}}})
}

func (m *MetaService) AddViewer(isNew bool, pathname string, isNewByPath bool) *domain.DataViewer {
	meta := m.GetMeta()
	oldViewer := int64(0)
	oldVisited := int64(0)
	if meta != nil {
		oldViewer = meta.Viewer
		oldVisited = meta.Visited
	}

	newViewer := oldViewer + 1
	newVisited := oldVisited
	if isNew {
		newVisited += 1
	}
	if meta != nil {
		m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "_id", Value: meta.Id}), bson.D{{"$set", bson.D{{"viewer", newViewer}, {"visited", newVisited}}}})
	}
	if strings.Contains(pathname, "post") {
		m.updateViewerByPathname(strings.ReplaceAll(pathname, "post", ""), isNew)
	}
	dateViewer := &domain.DateViewer{
		Date:    time.Now().Format(time.DateOnly),
		Viewer:  newViewer,
		Visited: newVisited,
	}
	m.ViewerService.SaveOrUpdateViewer(dateViewer)
	m.appendVisit(isNewByPath, pathname)

	return &domain.DataViewer{Visited: newVisited, Viewer: newViewer}
}

func (v *MetaService) appendVisit(isNew bool, pathname string) (bool, error) {
	now := time.Now()
	nowFormat := now.Format(time.DateOnly)
	todayVisit, err := v.VisitRepo.FindOne(mongodb.NewLogicalDefaultArray(bson.D{{"date", nowFormat}, {"pathname", pathname}}))
	if err != nil {
		return false, err
	}
	if todayVisit != nil {
		visited := todayVisit.Visited
		if isNew {
			visited = todayVisit.Visited + 1
		}
		update := bson.D{{"$set", bson.D{{"viewer", todayVisit.Viewer + 1}, {"visited", visited}, {"lastVisitedTime", time.Now()}}}}
		return v.VisitRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "_id", Value: todayVisit.Id}), update)
	} else {
		lastDay := now.AddDate(0, 0, -1)
		lastDayFormat := lastDay.Format(time.DateOnly)
		lastDayVisit, err := v.VisitRepo.FindOne(mongodb.NewLogicalDefaultArray(bson.D{{"date", lastDayFormat}, {"pathname", pathname}}))
		if err != nil {
			return false, nil
		}
		visited := lastDayVisit.Visited
		if isNew {
			visited += 1
		}

		visit := &domain.Visit{
			Date:            nowFormat,
			Viewer:          lastDayVisit.Viewer + 1,
			Visited:         visited,
			Pathname:        pathname,
			LastVisitedTime: time.Now(),
		}
		saved, err := v.VisitRepo.Save(visit)
		if err != nil {
			return false, err
		}
		return saved != "", nil
	}
}

func (a *MetaService) updateViewerByPathname(pathname string, isNew bool) {
	article := a.getArticleByIdOrPathname(pathname)
	if article != nil {
		oldViewer := article.Visited
		oldVisited := article.Visited
		newViewer := oldViewer + 1
		newVisited := oldVisited
		if isNew {
			newVisited = oldVisited + 1
		}
		a.ArticleRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "_id", Value: article.Id}), bson.D{{"$set", bson.D{{"visited", newVisited}, {"viewer", newViewer}}}})
	}
}

func (a *MetaService) getArticleByIdOrPathname(idOrPathname string) *domain.Article {
	var article *domain.Article
	id, err := strconv.Atoi(idOrPathname)
	if err == nil {
		article = a.getArticleById(int64(id))
	} else {
		article = a.getArticleByPathname(idOrPathname)
	}
	return article
}

func (a *MetaService) getArticleById(id int64) *domain.Article {
	filter := mongodb.NewLogicalDefault(bson.E{Key: "id", Value: id})
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	article, err := a.ArticleRepo.FindOne(filter)
	if err != nil {
		return &domain.Article{}
	}
	return article
}

func (a *MetaService) getArticleByPathname(pathname string) *domain.Article {
	filter := mongodb.NewLogicalDefault(bson.E{Key: "pathname", Value: pathname})
	filter.AppendLogical(mongodb.NewLogicalOrDefaultArray(DeleteFilter))
	article, err := a.ArticleRepo.FindOne(filter)
	if err != nil {
		return &domain.Article{}
	}
	return article
}
