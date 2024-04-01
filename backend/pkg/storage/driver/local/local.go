package local

import (
	"cc.allio/fusion/internal/token"
	"cc.allio/fusion/pkg/storage"
	"cc.allio/fusion/pkg/storage/driver"
	"cc.allio/fusion/pkg/util"
	"context"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const perm = 0744

type Driver struct {
	Policy *storage.Policy
}

func NewLocalDriver(policy *storage.Policy) driver.Handler {
	return &Driver{policy}
}

func (l *Driver) Init() error {
	// TODO
	return nil
}

func (l *Driver) Sign(ctx context.Context, file *storage.FileHeader) (string, error) {
	key := uuid.NewString()
	return token.Encrypt(key), nil
}

func (l *Driver) Upload(ctx context.Context, file *storage.FileStream) error {
	// obtain specifies save path about file system absolute path
	dest := util.AbsolutePath(file.SavePath, l.Policy.BaseDir)
	// create dir if not exist
	dir := filepath.Dir(dest)
	if !util.ExistFile(dir) {
		err := os.Mkdir(dir, perm)
		if err != nil {
			slog.Error("Failed to create directory. ", "err", err, "dir", dir)
			return err
		}
	}

	openMode := os.O_CREATE | os.O_RDWR
	out, err := os.OpenFile(dest, openMode, perm)
	if err != nil {
		slog.Error("Failed to Open file. ", "err", err, "dir", dir)
		return err
	}
	defer out.Close()
	// write to file
	_, err = io.Copy(out, file)
	return err
}

func (l *Driver) Remove(ctx context.Context, files []string) ([]string, error) {
	failedFiles := make([]string, 0, len(files))
	var err error
	for _, file := range files {
		filePath := util.AbsolutePath(file, l.Policy.BaseDir)
		if util.ExistFile(filePath) {
			err = os.Remove(filePath)
			if err != nil {
				slog.Error("Failed to remove file. ", "err", err, "filepath", filePath)
				failedFiles = append(failedFiles, file)
			}
		}
	}
	return failedFiles, err
}

func (l *Driver) Download(ctx context.Context, path string) (storage.RSCloser, error) {
	dest := util.AbsolutePath(path, l.Policy.BaseDir)
	file, err := os.OpenFile(dest, os.O_RDWR, perm)
	if err != nil {
		slog.Error("Failed to Open file. ", "err", err, "dest", dest)
		return nil, err
	}
	return file, err
}

func (l *Driver) Thumb(ctx context.Context, file *storage.FileHeader) (*storage.ContentResponse, error) {
	thumbFile := file.ThumbFile()

	rsCloser, err := l.Download(ctx, thumbFile)
	if err != nil {
		slog.Error("Failed to Download Thumb file. ", "err", err, "thumb", thumbFile)
		return nil, err
	}
	return &storage.ContentResponse{Redirect: false, Content: rsCloser}, nil
}

func (l *Driver) List(ctx context.Context, path string, recursive bool) ([]storage.Object, error) {
	var obj []storage.Object

	root := util.AbsolutePath(filepath.FromSlash(path), l.Policy.BaseDir)

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		// skip root path
		if path == root {
			return nil
		}

		if err != nil {
			slog.Error("Failed to walk folder", "err", err, "path", path)
			return filepath.SkipDir
		}

		// to relative path
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		obj = append(obj, storage.Object{
			Name:         info.Name(),
			RelativePath: filepath.ToSlash(rel),
			Source:       path,
			Size:         uint64(info.Size()),
			IsDir:        info.IsDir(),
			LastModify:   info.ModTime(),
		})

		if !recursive && info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})
	return obj, err
}
