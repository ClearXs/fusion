package filesystem

import (
	"cc.allio/fusion/pkg/storage"
	"cc.allio/fusion/pkg/storage/driver"
	"cc.allio/fusion/pkg/storage/driver/local"
	"errors"
)

type FileSystem struct {
	Policy  *storage.Policy
	Handler driver.Handler
}

func NewFs(policy *storage.Policy) *FileSystem {
	return &FileSystem{Policy: policy}
}

// LoadHandler by storage.Policy
func (fs *FileSystem) LoadHandler() error {
	var handler driver.Handler
	switch fs.Policy.Mode {
	case storage.LocalMode:
		handler = local.NewLocalDriver(fs.Policy)
	}
	if handler == nil {
		return errors.New("not found any policy load handler")
	}
	// init
	if err := handler.Init(); err != nil {
		return err
	}
	fs.Handler = handler
	return nil
}
