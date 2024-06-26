package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/pkg/mongodb"
	"cc.allio/fusion/pkg/storage"
	"cc.allio/fusion/pkg/util"
	"errors"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"
)

const (
	JwtSettingType    = "jwt"
	MenuSettingType   = "menu"
	StaticSettingType = "static"
	IsrSettingType    = "isr"
	LoginSettingType  = "login"
	HttpsSettingType  = "https"
	WalineSettingType = "waline"
	LayoutSettingType = "layout"
)

type SettingService struct {
	Cfg         *config.Config
	SettingRepo *repo.SettingsRepository
}

var SettingServiceSet = wire.NewSet(wire.Struct(new(SettingService), "*"))

// ---------------------- static ----------------------

func (s *SettingService) FindStaticSetting() *domain.StaticSetting {
	setting, err := s.SettingRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: StaticSettingType}))
	if err != nil {
		slog.Error("Find static setting has error", "err", err)
		return &domain.StaticSetting{}
	}
	value := setting.Value
	staticSetting := &domain.StaticSetting{
		Mode:            util.GetValue[string](value, "mode", storage.LocalMode),
		Endpoint:        util.GetValue[string](value, "endpoint", ""),
		AccessKeyID:     util.GetValue[string](value, "accessKeyID", ""),
		SecretAccessKey: util.GetValue[string](value, "secretAccessKey", ""),
		Bucket:          util.GetValue[string](value, "bucket", ""),
		BaseDir:         util.GetValue[string](value, "baseDir", ""),
		WaterMarkText:   util.GetValue[string](value, "waterMarkText", ""),
		EnableWaterMark: util.GetValue[bool](value, "enableWaterMark", false),
	}
	return staticSetting
}

func (s *SettingService) SaveOrUpdateStaticSetting(static *domain.StaticSetting) (bool, error) {
	outdated := s.FindStaticSetting()
	composite := util.Composition[domain.StaticSetting](outdated, static)
	value := util.EntityToMap[*domain.StaticSetting](composite)
	if outdated == nil {
		saved, err := s.SettingRepo.Save(&domain.Setting{Type: StaticSettingType, Value: value})
		if err != nil {
			return false, err
		}
		return saved > 0, nil
	} else {
		return s.SettingRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: StaticSettingType}), bson.D{{"$set", bson.D{{"value", value}}}})
	}
}

// ---------------------- login ----------------------

// FindLoginSetting if login setting is null, will be invoked SaveOrUpdateLoginSetting domain.DefaultLoginSetting
func (s *SettingService) FindLoginSetting() *domain.LoginSetting {
	setting, err := util.TryThen[domain.Setting](
		func() (*domain.Setting, error) {
			return s.SettingRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: LoginSettingType}))
		},
		func() (*domain.Setting, error) {
			value := util.EntityToMap[*domain.LoginSetting](&domain.DefaultLoginSetting)
			_, err := s.SettingRepo.Save(&domain.Setting{Type: LoginSettingType, Value: value})
			if err != nil {
				slog.Error("Failed to save default login setting.", "err", err)
				return nil, err
			}
			// re find
			return s.SettingRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: LoginSettingType}))
		},
		func(err error) bool {
			return errors.Is(err, mongo.ErrNoDocuments)
		},
	)
	if err != nil {
		slog.Error("Find https setting has error", "err", err)
		return &domain.LoginSetting{}
	}
	value := setting.Value
	loginSetting := &domain.LoginSetting{
		EnableMaxLoginRetry: value["enableMaxLoginRetry"].(bool),
		MaxRetryTimes:       int64(value["maxRetryTimes"].(float64)),
		DurationSeconds:     int64(value["durationSeconds"].(float64)),
		ExpiresIn:           int64(value["expiresIn"].(float64)),
	}
	return loginSetting
}

func (s *SettingService) SaveOrUpdateLoginSetting(layout *domain.LoginSetting) (bool, error) {
	outdated := s.FindLoginSetting()
	composite := util.Composition[domain.LoginSetting](outdated, layout)
	value := util.EntityToMap[*domain.LoginSetting](composite)
	if outdated == nil {
		saved, err := s.SettingRepo.Save(&domain.Setting{Type: LoginSettingType, Value: value})
		if err != nil {
			return false, err
		}
		return saved > 0, nil
	} else {
		return s.SettingRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: LoginSettingType}), bson.D{{"$set", bson.D{{"value", value}}}})
	}
}

// ---------------------- https ----------------------

func (s *SettingService) FindHttpsSetting() *domain.HttpSetting {
	setting, err := s.SettingRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: HttpsSettingType}))
	if err != nil {
		slog.Error("Find https setting has error", "err", err)
		return &domain.HttpSetting{}
	}
	value := setting.Value
	httpsSetting := &domain.HttpSetting{
		Redirect: value["redirect"].(bool),
	}
	return httpsSetting
}

// SaveOrUpdateHttpsSetting update https setting
func (s *SettingService) SaveOrUpdateHttpsSetting(credential *credential.HttpsSettingCredential) (bool, error) {
	https := s.FindHttpsSetting()
	if https != nil {
		filter := mongodb.NewLogicalDefault(bson.E{Key: "type", Value: HttpsSettingType})
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "value.redirect", Value: credential.Redirect}}}}
		successed, err := s.SettingRepo.Update(filter, update)
		if err != nil {
			slog.Error("update https setting has error", "err", err)
			return false, err
		}
		return successed, nil
	} else {
		httpSettings := &domain.Setting{Type: HttpsSettingType, Value: bson.M{"redirect": credential.Redirect}}
		_, err := s.SettingRepo.Save(httpSettings)
		if err != nil {
			slog.Error("save https setting has error", "err", err)
			return false, err
		}
		return true, nil
	}
}

// ---------------------- menu ----------------------

func (s *SettingService) FindMenuSettings() ([]*domain.MenuItem, error) {
	setting, err := s.SettingRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: MenuSettingType}))
	if err != nil {
		slog.Error("Find https setting has error", "err", err)
		return make([]*domain.MenuItem, 0), nil
	}
	value := setting.Value
	data, ok := value["value"].(map[string]interface{})
	if !ok {
		return nil, errors.New("not found any menus")
	}
	menuValues, ok := data["data"].(primitive.A)
	if !ok {
		return nil, errors.New("not found any menus")
	}

	menus := make([]*domain.MenuItem, 0)
	for _, dataValue := range menuValues {
		menu := util.MapToEntity[*domain.MenuItem](dataValue.(map[string]interface{}), &domain.MenuItem{})
		menus = append(menus, menu)
	}
	return menus, nil
}

func (s *SettingService) SaveOrUpdateMenuSettings(menus []*domain.MenuItem) (bool, error) {
	_, err := s.FindMenuSettings()
	value := make(map[string]any)
	value["data"] = util.EntityArrayToMapArray[*domain.MenuItem](menus)
	if err != nil {
		saved, err := s.SettingRepo.Save(&domain.Setting{Type: MenuSettingType, Value: value})
		if err != nil {
			return false, err
		}
		return saved > 0, nil
	} else {
		settings := &domain.Setting{Type: MenuSettingType, Value: value}
		return s.SettingRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: MenuSettingType}), bson.D{{"$set", bson.D{{"value", settings}}}})
	}
}

// ---------------------- waline ----------------------

func (s *SettingService) FindWalineSetting() *domain.WalineSetting {
	setting, err := s.SettingRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: WalineSettingType}))
	if err != nil {
		slog.Error("Find waline setting has error", "err", err)
		return nil
	}
	return &domain.WalineSetting{
		Email:             setting.Value["email"].(string),
		Enable:            setting.Value["enable"].(bool),
		ForceLoginComment: setting.Value["ForceLoginComment"].(bool),
	}
}

func (s *SettingService) SaveOrUpdateWalineSetting(waline *domain.WalineSetting) (bool, error) {
	outdated := s.FindWalineSetting()
	composite := util.Composition[domain.WalineSetting](outdated, waline)
	value := util.EntityToMap[*domain.WalineSetting](composite)
	if outdated == nil {
		saved, err := s.SettingRepo.Save(&domain.Setting{Type: WalineSettingType, Value: value})
		if err != nil {
			return false, err
		}
		return saved > 0, nil
	} else {
		return s.SettingRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: WalineSettingType}), bson.D{{"$set", bson.D{{"value", value}}}})
	}
}

// ---------------------- layou ----------------------

func (s *SettingService) FindLayoutSetting() *domain.LayoutSetting {
	setting, err := s.SettingRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: LayoutSettingType}))
	if err != nil {
		slog.Error("Find layout setting has error", "err", err)
		return nil
	}
	return &domain.LayoutSetting{
		Script: setting.Value["email"].(string),
		Html:   setting.Value["html"].(string),
		Css:    setting.Value["css"].(string),
		Head:   setting.Value["head"].(string),
	}
}

func (s *SettingService) SaveOrUpdateLayoutSetting(layout *domain.LayoutSetting) (bool, error) {
	outdated := s.FindLayoutSetting()
	composite := util.Composition[domain.LayoutSetting](outdated, layout)
	value := util.EntityToMap[*domain.LayoutSetting](composite)
	if outdated == nil {
		saved, err := s.SettingRepo.Save(&domain.Setting{Type: LayoutSettingType, Value: value})
		if err != nil {
			return false, err
		}
		return saved > 0, nil
	} else {
		return s.SettingRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: LayoutSettingType}), bson.D{{"$set", bson.D{{"value", value}}}})
	}
}

func (s *SettingService) FindIsrSetting() *domain.IsrSetting {
	setting, err := s.SettingRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: IsrSettingType}))
	if err != nil && setting != nil {
		slog.Error("Find layout setting has error", "err", err)
		return nil
	}
	return &domain.IsrSetting{
		Mode: setting.Value["mode"].(string),
	}
}

func (s *SettingService) SaveOrUpdateIsrSetting(isr *domain.IsrSetting) (bool, error) {
	outdated := s.FindIsrSetting()
	composite := util.Composition[domain.IsrSetting](outdated, isr)
	value := util.EntityToMap[*domain.IsrSetting](composite)
	if outdated == nil {
		saved, err := s.SettingRepo.Save(&domain.Setting{Type: IsrSettingType, Value: value})
		if err != nil {
			return false, err
		}
		return saved > 0, nil
	} else {
		return s.SettingRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: IsrSettingType}), bson.D{{"$set", bson.D{{"value", value}}}})
	}
}
