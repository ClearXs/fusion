package domain

type Setting struct {
	Id    int64          `json:"id" bson:"id"`
	Type  string         `bson:"type" json:"type"`
	Value map[string]any `bson:"value" json:"value"`
}

type StaticSetting struct {
	StorageType     string      `json:"storageType"`
	PicgoConfig     interface{} `json:"picgoConfig"`
	PicgoPlugins    interface{} `json:"picgoPlugins"`
	EnableWaterMark bool        `json:"enableWaterMark"`
	WaterMarkText   string      `json:"waterMarkText"`
	EnableWebp      bool        `json:"enableWebp"`
}

type LoginSetting struct {
	EnableMaxLoginRetry bool  `json:"enableMaxLoginRetry"`
	MaxRetryTimes       int64 `json:"maxRetryTimes"`
	DurationSeconds     int64 `json:"durationSeconds"`
	ExpiresIn           int64 `json:"expiresIn"`
}

type HttpSetting struct {
	Redirect bool `json:"redirect"`
}

type WalineSetting struct {
	Email             string `json:"email"`
	Enable            bool   `json:"smtp.enabled"`
	ForceLoginComment bool   `json:"forceLoginComment"`
}

type LayoutSetting struct {
	Script string `json:"script"`
	Html   string `json:"html"`
	Css    string `json:"css"`
	Head   string `json:"head"`
}
