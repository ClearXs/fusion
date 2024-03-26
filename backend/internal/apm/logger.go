package apm

import (
	bytes2 "bytes"
	"errors"
	"fmt"
	"github.com/asaskevich/EventBus"
	"github.com/goccy/go-json"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/h2non/gentleman.v2"
	"net/http"
	"os"
	"time"
)

// log type contains system business

type LogType = string

const (
	BusinessLogType LogType = "business"
	SystemLogType   LogType = "system"
)

// impl OpenObserve
// ref: https://openobserve.ai/docs/

type Logger struct {
	bus    EventBus.Bus
	config *LoggerConfig
	client *gentleman.Client
}

// Write combine apm logger and std output
func (l *Logger) Write(p []byte) (n int, err error) {
	// publish to apm
	// restructure record
	record := make(map[string]any)
	err = json.Unmarshal(p, &record)
	if err != nil {
		l.Put(string(p))
	} else {
		// add log type
		record["logType"] = SystemLogType
		l.Put(record)
	}
	// print to std output
	return os.Stdout.Write(p)
}

type LoggerConfig struct {
	// Enable to record log info
	Enable bool `yaml:"enable"`
	// Url for OpenObserve
	Url           string `yaml:"url"`
	Organization  string `yaml:"organization"`
	Stream        string `yaml:"stream"`
	Authorization string `yaml:"authorization"`
}

const collectTopic = "/fusion/apm/logger"

func NewLogger(bus EventBus.Bus, config *LoggerConfig) *Logger {
	return &Logger{bus: bus, config: config, client: gentleman.New()}
}

func Init(logger *Logger) {
	bus := logger.bus
	err := bus.SubscribeAsync(collectTopic, handlePut, true)
	if err != nil {
		fmt.Println("Failed to subscribe topic", "topic", collectTopic)
	}
}

// Record bson.D type log
func (l *Logger) Record(log bson.D) {
	l.Put(log)
}

// Put to APM System for OpenObserver base on event-bus
func (l *Logger) Put(v interface{}) {
	l.bus.Publish(collectTopic, l, v)
}

func handlePut(logger *Logger, v interface{}) {
	if !logger.config.Enable {
		// ignore...
		return
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		fmt.Println("Failed to serialization log", "err", err)
		return
	}
	client := logger.client
	config := logger.config
	// /api/{organization}/{stream}/_json
	path := fmt.Sprintf("/api/%s/%s/_json", config.Organization, config.Stream)
	_, err = client.Post().
		URL(config.Url).
		Path(path).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", config.Authorization).
		SetHeader("stream-name", config.Stream).
		Body(bytes2.NewBuffer(bytes)).
		Send()
	if err != nil {
		fmt.Println("Failed to send log. ", "err", err)
	}
}

type QueryResult struct {
	Data  []map[string]any `json:"hits"`
	Total int64            `json:"total"`
}

// Find query APM Logger system for OpenObserve
func (l *Logger) Find(search *Search) (Formatter[QueryResult], error) {
	config := l.config
	client := l.client

	// /api/{organization}/_search
	path := fmt.Sprintf("/api/%s/_search", config.Organization)
	requestData, err := json.Marshal(search)

	fmt.Println("Query log data for OpenObserver", "path", path, "request", string(requestData))
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
		fmt.Println("Failed query log data", "err", err)
		return nil, err
	}

	if !res.Ok || res.StatusCode != http.StatusOK {
		failedMessage := fmt.Sprintf("Failed to query log data for OpenObserver, path: [%s]  message: [%s], status code: [%d]", path, res.String(), res.StatusCode)
		fmt.Println(failedMessage)
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
