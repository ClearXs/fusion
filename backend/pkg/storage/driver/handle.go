package driver

import (
	"cc.allio/fusion/pkg/storage"
	"context"
)

// Handler 存储策略适配器
type Handler interface {

	// Init init storage handler
	Init() error

	// Upload file upload to dest path
	Upload(ctx context.Context, file *storage.FileStream) error

	// Remove ones or many path by files params.
	// returns failed file to string array and last error.
	Remove(ctx context.Context, files []string) ([]string, error)

	// Download get file content
	Download(ctx context.Context, path string) (storage.RSCloser, error)

	// Thumb obtain thumb about file
	Thumb(ctx context.Context, file *storage.FileHeader) (*storage.ContentResponse, error)

	// Sign to file
	Sign(ctx context.Context, file *storage.FileHeader) (string, error)

	// List specific path file list, recursive is true then recursion dir obtain file
	List(ctx context.Context, path string, recursive bool) ([]storage.Object, error)
}
