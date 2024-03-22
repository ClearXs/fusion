package storage

import "time"

type PolicyMode = string

const (
	LocalMode  = "local"
	MinioMode  = "minio"
	AliyunMode = "aliyun"
)

const osPathPrefix = "/fusion"

type Policy struct {
	Mode            PolicyMode `json:"mode"`
	Endpoint        string     `json:"endpoint"`
	AccessKeyID     string     `json:"accessKeyID"`
	SecretAccessKey string     `json:"secretAccessKey"`
	Bucket          string     `json:"bucket"`
	BaseDir         string     `json:"baseDir"`
}

// GenerateOsPath on FileHeader, the path will be os reality path
func (p *Policy) GenerateOsPath(f *FileHeader) string {
	path := f.FilePath
	date := time.Now().Format(time.DateOnly)
	return osPathPrefix + "/" + date + path
}
