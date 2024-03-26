package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/pkg/img"
	"cc.allio/fusion/pkg/storage"
	"cc.allio/fusion/pkg/storage/filesystem"
	"context"
	"github.com/google/wire"
	"golang.org/x/exp/slog"
	"time"
)

type FileService struct {
	Cfg              *config.Config
	StaticRepository *repo.StaticRepository
	SettingService   *SettingService
}

var FileServiceSet = wire.NewSet(wire.Struct(new(FileService), "*"))

func (f *FileService) Upload(ctx context.Context, header *storage.FileHeader) (*domain.Static, error) {
	policy := f.createPolicyBySetting()
	fs, err := f.chooseFs(policy)
	if err != nil {
		return nil, err
	}

	osPath := policy.GenerateOsPath(header)
	fileStream := &storage.FileStream{
		File:        header.File,
		Size:        header.Size,
		VirtualPath: header.FilePath,
		Name:        header.FilePath,
		SavePath:    osPath,
	}

	// sign to file
	sign, err := fs.Handler.Sign(ctx, header)
	if err != nil {
		return nil, err
	}

	if err = fs.Handler.Upload(ctx, fileStream); err != nil {
		return nil, err
	}

	// save file info to db
	var storageType domain.StorageType
	if img.IsImage(header.Filename) {
		storageType = domain.ImgStaticType
	} else {
		storageType = domain.FileCustomType
	}

	static := &domain.Static{
		Name:        header.Filename,
		StaticType:  policy.Mode,
		StorageType: storageType,
		FileType:    header.Ext,
		RealPath:    osPath,
		Sign:        sign,
		UpdatedAt:   time.Now(),
	}
	return static, err
}

// Delete by reality file path
func (f *FileService) Delete(ctx context.Context, filepath []string) error {
	policy := f.createPolicyBySetting()
	fs, err := f.chooseFs(policy)
	if err != nil {
		return err
	}
	slog.Info("prepare remove filepath", "filepath", filepath)
	failurePath, err := fs.Handler.Remove(ctx, filepath)
	if err != nil {
		slog.Error("Failed to remove file", "err", err, "failure path", failurePath)
		return err
	}
	return nil
}

func (f *FileService) createPolicyBySetting() *storage.Policy {
	staticSetting := f.SettingService.FindStaticSetting()
	return &storage.Policy{
		Mode:            staticSetting.Mode,
		Endpoint:        staticSetting.Endpoint,
		AccessKeyID:     staticSetting.AccessKeyID,
		SecretAccessKey: staticSetting.SecretAccessKey,
		Bucket:          staticSetting.Bucket,
		BaseDir:         staticSetting.BaseDir,
	}
}

func (f *FileService) chooseFs(policy *storage.Policy) (*filesystem.FileSystem, error) {
	newFs := filesystem.NewFs(policy)
	// load handler
	err := newFs.LoadHandler()
	if err != nil {
		return nil, err
	}
	return newFs, nil
}
