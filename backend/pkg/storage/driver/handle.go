package driver

import (
	"cc.allio/fusion/pkg/storage"
	"context"
)

// Handler 存储策略适配器
type Handler interface {

	// Upload file upload to dest path
	Upload(ctx context.Context, file *storage.FileStream) error

	// Remove ones or many path, and returns failed file to array string and second last error
	Remove(ctx context.Context, files []string) ([]string, error)

	// Download get file content
	Download(ctx context.Context, path string) (storage.RSCloser, error)

	// Thumb obtain thumb about file
	Thumb(ctx context.Context, file *storage.File) (*storage.ContentResponse, error)

	// List specific path file list, recursive is true then recursion dir obtain file
	List(ctx context.Context, path string, recursive bool) ([]storage.Object, error)
}
