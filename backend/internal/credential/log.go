package credential

type LogSearchOption struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	LogType  string `json:"eventType"`
}
