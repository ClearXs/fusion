package svr

import (
	"bytes"
	"cc.allio/fusion/config"
	"context"
	"github.com/goccy/go-json"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
	"os"
)

type CaddyService struct {
	Cfg             *config.Config
	SettingsService SettingService
}

var CaddyServiceSet = wire.NewSet(wire.Struct(new(CaddyService), "*"))

func Init(caddy *CaddyService) (func(), error) {
	https := caddy.SettingsService.FindHttpsSetting()
	caddy.SetRedirect(https.Redirect)
	return func() {
		caddy.clear()
	}, nil
}

func (caddy *CaddyService) clear() {
	// TODO
}

func (caddy *CaddyService) ClearLog() {
	var caddyLogPath = "/var/log/caddy.log"

	err := os.WriteFile(caddyLogPath, []byte{}, 0777)
	if err != nil {
		slog.Error("clear caddy log failed from path /var/log/caddy.log", "err", err)
	}
}

func (caddy *CaddyService) GetLog() (string, error) {
	var caddyLogPath = "/var/log/caddy.log"
	data, err := os.ReadFile(caddyLogPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (caddy *CaddyService) GetConfig() (map[string]any, error) {
	var url = "http://127.0.0.1:2019/config"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	binaries, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data := make(map[string]any)
	if err = json.Unmarshal(binaries, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (caddy *CaddyService) SetRedirect(redirect bool) bool {
	var (
		req *http.Request
		err error
		url = "http://127.0.0.1:2019/config/apps/http/servers/srv1/listener_wrappers"
	)
	if !redirect {
		req, err = http.NewRequest(http.MethodDelete, url, nil)
	} else {
		body := bson.D{{"wrapper", "http_redirect"}}
		data, err := bson.Marshal(body)
		if err != nil {
			return false
		}
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	}
	if err != nil {
		slog.Error("set redirect failed from remote svr", "url", url, "err", err)
		return false
	}
	defer func() {
		if req != nil {
			req.Clone(context.TODO())
		}
	}()
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("set redirect failed from remote svr", "url", url, "err", err)
		return false
	}
	slog.Info("open redirect success")
	return true
}
