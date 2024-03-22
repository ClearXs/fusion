package credential

import "cc.allio/fusion/internal/domain"

type StaticSearchOption struct {
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	StaticType domain.StaticType `json:"staticType"`
}

type FilePathCredential struct {
	Src   string `json:"src"`
	IsNew bool   `json:"isNew"`
}
