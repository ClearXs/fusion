package event

import (
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/env"
	"cc.allio/fusion/pkg/fetch"
	"cc.allio/fusion/pkg/guc"
	"github.com/asaskevich/EventBus"
	"github.com/google/wire"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc"
	"golang.org/x/exp/slog"
	"strconv"
	"strings"
	"time"
)

// implementation backend publish isr notification

type Path = string

type IsrEventBus struct {
	WebSiteUrl string
	Bus        EventBus.Bus
	Service    *svr.Service
}

const DefaultRetryCount = 3
const DelayTimeout = time.Duration(3000) // default 3 second

var IsrEventBusSet = wire.NewSet(NewIsrEventBus)

func NewIsrEventBus(Service *svr.Service) *IsrEventBus {
	return &IsrEventBus{WebSiteUrl: env.WebSiteUrl, Bus: EventBus.New(), Service: Service}
}

// handle fetch to path
func (isr *IsrEventBus) handle(path string, args ...interface{}) {
	slog.Info("trigger isr rendering", "path", path, "args", args)
	_, err := fetch.Get(isr.WebSiteUrl, fetch.D{{"path", path}})
	if err != nil {
		slog.Error("failed to trigger isr rendering", "path", path)
	}
}

// retryHandle
func (isr *IsrEventBus) retryHandle(retryCount int, path string, args ...interface{}) {
	slog.Info("retry trigger isr rendering", "path", path, "retryCount", retryCount, "delay", DefaultRetryCount, "args", args)
	var wg conc.WaitGroup
	wg.Go(func() {
		results := guc.Retry(retryCount, DelayTimeout, func() error {
			_, err := fetch.Get(isr.WebSiteUrl, fetch.D{{"path", path}})
			if err != nil {
				slog.Error("failed to trigger isr rendering", "path", path)
			}
			return nil
		})
		for r := range results {
			slog.Info("retry message", "path", path, "r", r)
		}
	})
	wg.Wait()
}

// Active publish path as topic into on Event-Bus.
// subscriber subscribe path through handle process
func (isr *IsrEventBus) Active(path Path, args ...interface{}) {
	slog.Info("active isr path", "path", path)
	bus := isr.Bus
	if bus.HasCallback(path) {
		err := bus.SubscribeAsync(path, isr.handle, true)
		if err != nil {
			slog.Error("subscribe path has error", "err", err, "path", path)
		}
	}
	bus.Publish(path, path, args)
}

// ActiveRetry same as Active but subscriber through retryHandle process
func (isr *IsrEventBus) ActiveRetry(path Path, args ...interface{}) {
	slog.Info("active isr path", "path", path)
	bus := isr.Bus
	if bus.HasCallback(path) {
		err := bus.SubscribeAsync(path, isr.retryHandle, true)
		if err != nil {
			slog.Error("subscribe path has error", "err", err, "path", path)
		}
	}
	bus.Publish(path, DefaultRetryCount, path, args)
}

func (isr *IsrEventBus) ActiveAll(args ...interface{}) {
	isr.ActiveCategory(args...)
	isr.ActivePost(args...)
	isr.ActivePage(args...)
	isr.ActiveTag(args...)
}

// ActiveCategory active system all category make it rendering
func (isr *IsrEventBus) ActiveCategory(args ...interface{}) {
	categoryKeys := isr.Service.CategoryService.GetAllCategoryKeys()
	for _, categoryKey := range categoryKeys {
		isr.Active("/category/"+isr.encodePath(categoryKey), args...)
	}
}

// ActivePost active system all post make it rendering
func (isr *IsrEventBus) ActivePost(args ...interface{}) {
	articles := isr.Service.ArticleService.GetAll("list", true, true)
	for _, article := range articles {
		if lo.IsNotEmpty(article.Pathname) {
			isr.Active("/post/"+article.Pathname, args...)
		} else {
			isr.Active("/post/"+strconv.FormatInt(article.Id, 10), args...)
		}
	}
}

// ActivePage active system first page make it rendering
func (isr *IsrEventBus) ActivePage(args ...interface{}) {
	totalNum := isr.Service.ArticleService.GetTotalNum(false)
	firstNum := totalNum / 5
	for i := int64(1); i <= firstNum; i++ {
		isr.Active("/page/"+strconv.FormatInt(i, 10), args...)
	}
}

// ActiveTag active system all tag make it rendering
func (isr *IsrEventBus) ActiveTag(args ...interface{}) {
	tags := isr.Service.TagService.GetAllTags(false)
	for _, tag := range tags {
		isr.Active("/tag/"+isr.encodePath(tag), args...)
	}
}

// ActiveAbout active about info
func (isr *IsrEventBus) ActiveAbout(args ...interface{}) {
	isr.ActiveRetry("/about", args...)
}

// ActiveLink active link
func (isr *IsrEventBus) ActiveLink(args ...interface{}) {
	isr.ActiveRetry("/link", args...)
}

func (isr *IsrEventBus) encodePath(path string) string {
	encode := strings.ReplaceAll(path, "#", "%23")
	encode = strings.ReplaceAll(encode, "/", "%2F")
	return encode
}
