package credential

import "time"

type MetaCredential struct {
	Version       string      `json:"version"`
	LatestVersion string      `json:"latestVersion"`
	UpdateAt      time.Time   `json:"updateAt"`
	User          interface{} `json:"user"`
	BaseUrl       string      `json:"baseUrl"`
	EnableComment string      `json:"enableComment"`
	AllowDomains  string      `json:"allowDomains"`
}
