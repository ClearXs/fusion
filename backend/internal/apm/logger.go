package apm

import (
	bytes2 "bytes"
	"errors"
	"fmt"
	"github.com/asaskevich/EventBus"
	"github.com/goccy/go-json"
	"github.com/google/wire"
	"golang.org/x/exp/slog"
	"gopkg.in/h2non/gentleman.v2"
	"net/http"
	"time"
)

// impl OpenObserve
// ref: https://openobserve.ai/docs/

type Logger struct {
	bus    EventBus.Bus
	config *Config
	client *gentleman.Client
}

type Config struct {
	// Url for OpenObserve
	Url           string `yaml:"url"`
	Organization  string `yaml:"organization"`
	Stream        string `yaml:"stream"`
	Authorization string `yaml:"authorization"`
}

const collectTopic = "/fusion/apm/logger"

func NewLogger(bus EventBus.Bus, config *Config) *Logger {
	return &Logger{bus: bus, config: config, client: gentleman.New()}
}

func Init(logger *Logger) {
	bus := logger.bus
	err := bus.SubscribeAsync(collectTopic, handlePut, true)
	if err != nil {
		slog.Error("Failed to subscribe topic", "topic", collectTopic)
	}
}

var LoggerSet = wire.NewSet(NewLogger, Init)

// Put to APM System for OpenObserver base on event-bus
func (l *Logger) Put(v interface{}) {
	l.bus.Publish(collectTopic, l, v)
}

func handlePut(logger *Logger, v interface{}) {
	slog.Debug("handle log send to OpenObserve", "log", v)
	bytes, err := json.Marshal(v)
	if err != nil {
		slog.Error("Failed to serialization log", "err", err)
		return
	}
	client := logger.client
	config := logger.config
	// /api/{organization}/{stream}/_json
	path := fmt.Sprintf("/api/%s/%s/_json", config.Organization, config.Stream)
	res, err := client.Post().
		URL(config.Url).
		Path(path).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", config.Authorization).
		SetHeader("stream-name", config.Stream).
		Body(bytes2.NewBuffer(bytes)).
		Send()
	if err != nil {
		slog.Error("Failed to send log. ", "err", err)
	}

	// record send information
	data := string(res.Bytes())
	slog.Debug("send log", "ok status", res.Ok, "data", data)
}

type QueryResult struct {
	Data  []map[string]any `json:"hits"`
	Total int64            `json:"total"`
}

// FindMap query APM Logger system for OpenObserve
func (l *Logger) FindMap(search *Search) (Formatter[QueryResult], error) {
	config := l.config
	client := l.client

	// /api/{organization}/_search
	path := fmt.Sprintf("/api/%s/_search", config.Organization)
	requestData, err := json.Marshal(search)

	slog.Debug("Query log data for OpenObserver", "path", path, "request", string(requestData))

	if err != nil {
		return nil, err
	}

	res, err := client.
		Post().
		URL(config.Url).
		Path(path).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", config.Authorization).
		Body(bytes2.NewBuffer(requestData)).
		Send()

	if err != nil {
		slog.Error("Failed query log data", "err", err)
		return nil, err
	}

	if !res.Ok || res.StatusCode != http.StatusOK {
		failedMessage := fmt.Sprintf("Failed to query log data for OpenObserver, path: [%s]  message: [%s], status code: [%d]", path, res.String(), res.StatusCode)
		slog.Error(failedMessage)
		return nil, errors.New(failedMessage)
	}

	return NewFormat[QueryResult](res.Bytes()), nil
}

type SQLMode = string

// define OpenObserve support to SQLMode
// context: default is context mode, you cann't use limit group by in query.sql, and will use the SQL result as a context of aggregations. aggregation will get result from context.
// full: in full mode, you can write a full SQL in query.sql, it supports limit group by and keywords, but it doesn't support aggregation.
const (
	ContextSQLMode SQLMode = "context"
	FullSQLMode    SQLMode = "full"
)

type Aggs = []string

type SearchQuery struct {
	SQL            string        `json:"sql"`
	StartTime      time.Duration `json:"start_time,omitempty"`
	EndTime        time.Duration `json:"end_time,omitempty"`
	From           int           `json:"from,omitempty"`
	Size           int           `json:"size,omitempty"`
	TrackTotalHits bool          `json:"track_total_hits,omitempty"`
	SQLMode        SQLMode       `json:"sql_mode"`
}

// Search OpenObserve query structure
// see https://openobserve.ai/docs/api/search/search/#request
type Search struct {
	SearchQuery `json:"query"`
	Aggs        `json:"aggs,omitempty"`
}

func NewSearch() *Search {
	return &Search{
		SearchQuery: SearchQuery{SQLMode: ContextSQLMode},
		Aggs:        nil,
	}
}

// build strategy

func (s *Search) Agg(agg string) *Search {
	if s.Aggs == nil {
		s.Aggs = make(Aggs, 0)
	}
	s.Aggs = append(s.Aggs, agg)
	return s
}

func (s *Search) SQL(sql string) *Search {
	s.SearchQuery.SQL = sql
	return s
}

func (s *Search) Interval(startTime time.Time, endTime time.Time) *Search {
	s.SearchQuery.StartTime = time.Duration(startTime.UnixMilli())
	s.SearchQuery.EndTime = time.Duration(endTime.UnixMilli())
	return s
}

func (s *Search) Page(page int, pageSize int) *Search {
	s.SearchQuery.From = (page - 1) * pageSize
	s.SearchQuery.Size = pageSize
	return s
}

func (s *Search) TrackTotalHits(trackTotalHits bool) *Search {
	s.SearchQuery.TrackTotalHits = trackTotalHits
	return s
}

func (s *Search) SQLMode(mode SQLMode) *Search {
	s.SearchQuery.SQL = mode
	return s
}
