package domain

import (
	"time"
)

type Viewer struct {
	Id        int64     `json:"id" bson:"id"`
	Visited   int64     `json:"visited" bson:"visited"`
	Viewer    int64     `json:"viewer" bson:"viewer"`
	Date      string    `json:"date" bson:"date"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type ViewerGrid struct {
	Grid GridDateView `json:"grid"`
	Add  DataViewer   `json:"add"`
	Now  DataViewer   `json:"now"`
}

type DataViewer struct {
	Viewer  int64 `json:"viewer"`
	Visited int64 `json:"visited"`
}

type GridDateView struct {
	Total []DateViewer `json:"total"`
	Each  []DateViewer `json:"each"`
}

type DateViewer struct {
	Date    string `json:"date"`
	Viewer  int64  `json:"viewer"`
	Visited int64  `json:"visited"`
}

type OverviewTabData struct {
	Total struct {
		WordCount  int64 `json:"wordCount"`
		ArticleNum int64 `json:"articleNum"`
	} `json:"total"`

	Viewer ViewerGrid `json:"viewer"`
	Link   struct {
		BaseUrl       string      `json:"baseUrl"`
		EnableComment TrueOrFalse `json:"enableComment"`
	} `json:"link"`
}

type ViewerTabData struct {
	EnableGA                bool       `json:"enableGA"`
	EnableBaidu             bool       `json:"enableBaidu"`
	TopViewer               []*Article `json:"topViewer" `
	TopVisited              []*Article `json:"topVisited" `
	RecentVisitArticles     []*Article `json:"recentVisitArticles"`
	SiteLastVisitedTime     time.Time  `json:"siteLastVisitedTime" `
	SiteLastVisitedPathname string     `json:"siteLastVisitedPathname"`
	TotalViewer             int64      `json:"totalViewer" `
	TotalVisited            int64      `json:"totalVisited" `
	MaxArticleViewer        int64      `json:"maxArticleViewer"`
	MaxArticleVisited       int64      `json:"maxArticleVisited"`
}

type ArticleTabData struct {
	ArticleNum      int64               `json:"articleNum" `
	CategoryNum     int64               `json:"categoryNum" `
	TagNum          int64               `json:"tagNum" `
	WordNum         int64               `json:"wordNum" `
	CategoryPieData []*TypeValue[int64] `json:"categoryPieData"`
	ColumnData      []*TypeValue[int64] `json:"columnData" `
}
