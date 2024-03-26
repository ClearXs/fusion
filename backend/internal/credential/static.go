package credential

type StaticSearchOption struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	StaticType string `json:"staticType"`
}

type FilePathCredential struct {
	Src   string `json:"src"`
	IsNew bool   `json:"isNew"`
}
