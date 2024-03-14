package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/pkg/mongodb"
	"cc.allio/fusion/pkg/utils"
	"errors"
	"github.com/google/wire"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slog"
	"time"
)

type MetaService struct {
	Cfg      *config.Config
	MetaRepo *repo.MetaRepository
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
	if meta != nil {
		return false, errors.New("meta is empty")
	}
	value := utils.Composition[domain.SiteInfo](meta.SiteInfo, siteInfo)

	elements := utils.ToBsonElements(value)
	return m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: meta.Id}), bson.D{{"siteInfo", elements}})
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
	return m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: meta.Id}), bson.D{{"links", links}})
}

func (m *MetaService) DeleteLinkByName(name string) (bool, error) {
	meta := m.GetMeta()
	if meta.Links == nil {
		return false, nil
	}
	links := meta.Links
	links = lo.Filter(links, func(item *domain.LinkItem, index int) bool { return item.Name == name })
	return m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: meta.Id}), bson.D{{"links", links}})
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
	return m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: meta.Id}), bson.D{{"rewards", rewards}})
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
	return m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: meta.Id}), bson.D{{"rewards", rewards}})
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
	return m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: meta.Id}), bson.D{{"socials", socials}})
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
	return m.MetaRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "id", Value: meta.Id}), bson.D{{"socials", socials}})
}
