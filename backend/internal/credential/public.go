package credential

import (
	"cc.allio/fusion/internal/domain"
	"time"
)

type PublicArticle struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	Top       int64     `json:"top"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SiteMeta struct {
	Version string   `json:"version"`
	Tags    []string `json:"tags"`
	Meta    struct {
		Id             string               `json:"id"`
		Links          []*domain.LinkItem   `json:"links"`
		Socials        []*domain.SocialItem `json:"socials"`
		Menus          []*domain.MenuItem   `json:"menus"`
		Rewards        []*domain.RewardItem `json:"rewards"`
		About          *domain.About        `json:"about"`
		SiteInfo       *domain.SiteInfo     `json:"siteInfo"`
		Viewer         int64                `json:"viewer"`
		Visited        int64                `json:"visited"`
		TotalWordCount int64                `json:"totalWordCount"`
		Categories     []*domain.Category   `json:"categories"`
	} `json:"meta"`
	Menus         []*domain.MenuItem    `json:"menus"`
	TotalArticles int64                 `json:"totalArticles"`
	TotalWords    int64                 `json:"totalWords"`
	LayoutSetting *domain.LayoutSetting `json:"layoutSetting"`
}
