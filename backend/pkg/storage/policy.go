package storage

type PolicyMode = string

const (
	LocalMode  = "local"
	MinioMode  = "minio"
	AliyunMode = "aliyun"
)

type Policy struct {
	Mode            PolicyMode `json:"mode"`
	Endpoint        string     `json:"endpoint"`
	AccessKeyID     string     `json:"accessKeyID"`
	SecretAccessKey string     `json:"secretAccessKey"`
	Bucket          string     `json:"bucket"`
	BaseDir         string     `json:"baseDir"`
}
