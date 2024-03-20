package storage

type File struct {
	Name string
	Path string
}

// ThumbFile returns thumb file name
func (f *File) ThumbFile() string {
	return f.Path + "._thumb"
}
