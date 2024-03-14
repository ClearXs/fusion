package credential

import "cc.allio/fusion/internal/domain"

type DraftSearchOptionCredential struct {
	Page          int              `json:"page"`
	PageSize      int              `json:"pageSize"`
	ToListView    bool             `json:"toListView"`
	Category      string           `json:"category"`
	Tags          string           `json:"tags"`
	Title         string           `json:"title"`
	SortCreatedAt domain.SortOrder `json:"sortCreatedAt"`
	StartTime     string           `json:"startTime"`
	EndTime       string           `json:"endTime"`
}

type DraftPublishCredential struct {
	Hidden    bool   `json:"hidden"`
	Pathname  string `json:"pathname"`
	Private   bool   `json:"private"`
	Password  string `json:"password"`
	Copyright string `json:"copyright"`
}
