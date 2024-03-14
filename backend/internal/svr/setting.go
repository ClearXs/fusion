package svr

import (
	"cc.allio/fusion/config"
	"cc.allio/fusion/internal/credential"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/repo"
	"cc.allio/fusion/pkg/mongodb"
	"cc.allio/fusion/pkg/utils"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
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
		StorageType:     value["storageType"].(string),
		PicgoConfig:     value["picgoConfig"],
		PicgoPlugins:    value["picgoPlugins"],
		EnableWaterMark: value["enableWaterMark"].(bool),
		WaterMarkText:   value["WaterMarkText"].(string),
		EnableWebp:      value["enableWebp"].(bool),
	}
	return staticSetting
}

func (s *SettingService) SaveOrUpdateStaticSetting(static *domain.StaticSetting) (bool, error) {
	outdated := s.FindStaticSetting()
	composite := utils.Composition[domain.StaticSetting](outdated, static)
	value := utils.EntityToMap[*domain.StaticSetting](composite)
	if outdated == nil {
		saved, err := s.SettingRepo.Save(&domain.Setting{Type: StaticSettingType, Value: value})
		if err != nil {
			return false, err
		}
		return saved > 0, nil
	} else {
		return s.SettingRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: StaticSettingType}), bson.D{{"value", value}})
	}
}

// ---------------------- login ----------------------

func (s *SettingService) FindLoginSetting() *domain.LoginSetting {
	setting, err := s.SettingRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: LoginSettingType}))
	if err != nil {
		slog.Error("Find static setting has error", "err", err)
		return &domain.LoginSetting{}
	}
	value := setting.Value
	loginSetting := &domain.LoginSetting{
		EnableMaxLoginRetry: value["enableMaxLoginRetry"].(bool),
		MaxRetryTimes:       value["maxRetryTimes"].(int64),
		DurationSeconds:     value["durationSeconds"].(int64),
		ExpiresIn:           value["expiresIn"].(int64),
	}
	return loginSetting
}

func (s *SettingService) SaveOrUpdateLoginSetting(layout *domain.LoginSetting) (bool, error) {
	outdated := s.FindLoginSetting()
	composite := utils.Composition[domain.LoginSetting](outdated, layout)
	value := utils.EntityToMap[*domain.LoginSetting](composite)
	if outdated == nil {
		saved, err := s.SettingRepo.Save(&domain.Setting{Type: LoginSettingType, Value: value})
		if err != nil {
			return false, err
		}
		return saved > 0, nil
	} else {
		return s.SettingRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: LoginSettingType}), bson.D{{"value", value}})
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

func (s *SettingService) FindMenuSettings() []*domain.MenuItem {
	setting, err := s.SettingRepo.FindOne(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: MenuSettingType}))
	if err != nil {
		slog.Error("Find https setting has error", "err", err)
		return make([]*domain.MenuItem, 0)
	}
	value := setting.Value
	data := value["data"].([]map[string]any)
	menus := utils.MapArrayToEntityArray[*domain.MenuItem](data, func() *domain.MenuItem { return &domain.MenuItem{} })
	return menus
}

func (s *SettingService) SaveOrUpdateMenuSettings(menu *domain.MenuItem) (bool, error) {
	menus := s.FindMenuSettings()
	value := make(map[string]any)
	value["data"] = utils.EntityArrayToMapArray[*domain.MenuItem](append(menus, menu))
	if len(menus) == 0 {
		saved, err := s.SettingRepo.Save(&domain.Setting{Type: MenuSettingType, Value: value})
		if err != nil {
			return false, err
		}
		return saved > 0, nil
	} else {
		settings := &domain.Setting{Type: MenuSettingType, Value: value}
		return s.SettingRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: MenuSettingType}), bson.D{{"value", settings}})
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
	composite := utils.Composition[domain.WalineSetting](outdated, waline)
	value := utils.EntityToMap[*domain.WalineSetting](composite)
	if outdated == nil {
		saved, err := s.SettingRepo.Save(&domain.Setting{Type: WalineSettingType, Value: value})
		if err != nil {
			return false, err
		}
		return saved > 0, nil
	} else {
		return s.SettingRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: WalineSettingType}), bson.D{{"value", value}})
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
	composite := utils.Composition[domain.LayoutSetting](outdated, layout)
	value := utils.EntityToMap[*domain.LayoutSetting](composite)
	if outdated == nil {
		saved, err := s.SettingRepo.Save(&domain.Setting{Type: LayoutSettingType, Value: value})
		if err != nil {
			return false, err
		}
		return saved > 0, nil
	} else {
		return s.SettingRepo.Update(mongodb.NewLogicalDefault(bson.E{Key: "type", Value: LayoutSettingType}), bson.D{{"value", value}})
	}
}
