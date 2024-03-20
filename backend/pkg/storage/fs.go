package storage

import (
	"io"
	"time"
)

// FileStream describe file info
type FileStream struct {
	File        io.ReadCloser
	Seeker      io.Seeker
	Size        uint64
	VirtualPath string
	Name        string
	SavePath    string
}

func (f *FileStream) Read(p []byte) (n int, err error) {
	if f.File != nil {
		return f.File.Read(p)
	}
	return 0, io.EOF
}

type RSCloser interface {
	io.ReadSeeker
	io.Closer
}

// ContentResponse file content
type ContentResponse struct {
	Redirect bool
	Content  RSCloser
	URL      string
	MaxAge   int
}

// Object list file object
type Object struct {
	Name         string    `json:"name"`
	RelativePath string    `json:"relative_path"`
	Source       string    `json:"source"`
	Size         uint64    `json:"size"`
	IsDir        bool      `json:"is_dir"`
	LastModify   time.Time `json:"last_modify"`
}
