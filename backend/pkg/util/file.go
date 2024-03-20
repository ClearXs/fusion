package util

import (
	"golang.org/x/exp/slog"
	"os"
	"path/filepath"
)

// ExistFile reports whether the name file or dir exist
func ExistFile(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// AbsolutePath obtain relative path, if path not absolute, return project root append path
func AbsolutePath(path string, baseDir string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(baseDir, path)
}

// CreateNestedFile specifies file path nested create direction finally create file
func CreateNestedFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	if !ExistFile(dir) {
		err := os.MkdirAll(path, 0700)
		if err != nil {
			slog.Error("Failed to create nested file. ", "err", err, "path", path)
			return nil, err
		}
	}
	return os.Create(path)
}
