package fetch

import (
	"github.com/goccy/go-json"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
)

type P struct {
	Key   string
	Value string
}

type D []P

// Get build Http 'Get' request by specific url and params
func Get(url string, params D) (map[string]any, error) {
	slog.Info("Build http 'Get' request", "url", url, "params", params)

	// append url a pair of param key and value
	if params != nil {
		paths := lo.Reduce[P, string](params, func(agg string, item P, index int) string {
			pair := item.Key + "=" + item.Value
			if lo.IsEmpty(agg) {
				return pair
			} else {
				return agg + "&" + pair
			}
		}, "")

		if lo.IsNotEmpty(paths) {
			url = url + "?" + paths
		}
	}

	// build new request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		slog.Error("failed to build http 'Get'", "err", err)
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Http 'Get' request process has error", "err", err, "url", url, "params", params)
		return nil, err
	}
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	// read and handle data to map
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
