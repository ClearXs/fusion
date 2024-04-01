package domain

import "time"

type Setting struct {
	Id    int64          `json:"id" bson:"id"`
	Type  string         `bson:"type" json:"type"`
	Value map[string]any `bson:"value" json:"value"`
}

type StaticSetting struct {
	Mode            string `json:"mode" bson:"mode"`
	Endpoint        string `json:"endpoint" bson:"endpoint"`
	AccessKeyID     string `json:"accessKeyID" bson:"accessKeyID"`
	SecretAccessKey string `json:"secretAccessKey" bson:"secretAccessKey"`
	Bucket          string `json:"bucket" bson:"bucket"`
	BaseDir         string `json:"baseDir" bson:"baseDir"`
	WaterMarkText   string `json:"waterMarkText" bson:"waterMarkText"`
	EnableWaterMark bool   `json:"enableWaterMark" bson:"enableWaterMark"`
}

type LoginSetting struct {
	EnableMaxLoginRetry bool  `json:"enableMaxLoginRetry"`
	MaxRetryTimes       int64 `json:"maxRetryTimes"`
	DurationSeconds     int64 `json:"durationSeconds"`
	ExpiresIn           int64 `json:"expiresIn"`
}

var DefaultLoginSetting = LoginSetting{
	EnableMaxLoginRetry: false,
	MaxRetryTimes:       int64(time.Minute),
	DurationSeconds:     int64(30 * time.Second),
	ExpiresIn:           int64(10 * time.Hour),
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

type IsrSetting struct {
	Mode string `json:"mode"`
}
