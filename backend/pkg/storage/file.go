package storage

import "io"

type FileHeader struct {
	FilePath string
	Filename string
	Header   map[string][]string
	Size     uint64
	File     io.ReadCloser
	Ext      string
}

// ThumbFile returns thumb file name
func (f *FileHeader) ThumbFile() string {
	return f.FilePath + "._thumb"
}
