package credential

import "cc.allio/fusion/internal/domain"

type ArticleSearchOptionCredential struct {
	Page          int              `json:"page"`
	PageSize      int              `json:"pageSize"`
	RegMatch      bool             `json:"regMatch"`
	Category      string           `json:"category"`
	Tags          string           `json:"tags"`
	Title         string           `json:"title"`
	SortCreatedAt domain.SortOrder `json:"sortCreatedAt"`
	SortTop       domain.SortOrder `json:"sortTop"`
	StartTime     string           `json:"startTime"`
	EndTime       string           `json:"endTime"`
	SortViewer    string           `json:"sortViewer"`
	ToListView    bool             `json:"toListView"`
	WithWordCount bool             `json:"withWordCount"`
	Author        string           `json:"author"`
}

type ArticleSearchLinkCredential struct {
	Link string `json:"link"`
}

type ArticlePasswordCredential struct {
	Password string `json:"password"`
}

type AlternateArticle struct {
	Pre     *domain.Article `json:"pre"`
	Article *domain.Article `json:"article"`
	Next    *domain.Article `json:"next"`
}
