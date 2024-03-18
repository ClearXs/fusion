package event

import (
	"bytes"
	"cc.allio/fusion/internal/domain"
	"cc.allio/fusion/internal/svr"
	"cc.allio/fusion/pkg/guc"
	"github.com/goccy/go-json"
	"github.com/google/wire"
	"github.com/robertkrimen/otto"
	"github.com/sourcegraph/conc"
	"golang.org/x/exp/slog"
	"html/template"
	"time"
)

type ScriptEngine struct {
	Timeout         time.Duration
	PipelineService *svr.PipelineService
	tmpl            *template.Template
}

var ScriptEngineSet = wire.NewSet(NewScripEngine)

type CodeResultStatus = string

const (
	SuccessStatus CodeResultStatus = "success"
	ErrorStatus   CodeResultStatus = "error"
)

type CodeResult struct {
	Logs   []string         `json:"logs"`
	Output interface{}      `json:"output"`
	Status CodeResultStatus `json:"status"`
}

const ScriptTemplate = `
  try {
    {{.}}
    send(input, "success"); // send success
  } catch(err) {
    input.error = err;
    send(input, "error"); // send success
  }
`

func NewScripEngine(PipelineService *svr.PipelineService) *ScriptEngine {
	tmpl, err := template.New("script").Parse(ScriptTemplate)
	if err != nil {
		panic(err)
	}
	return &ScriptEngine{Timeout: guc.DelayTimeout, PipelineService: PipelineService, tmpl: tmpl}
}

// Dispatch execute pipeline script
func (s *ScriptEngine) Dispatch(pipeline *domain.Pipeline, args ...interface{}) <-chan *CodeResult {
	r := make(chan *CodeResult)
	// build script text
	buf := new(bytes.Buffer)
	err := s.tmpl.Execute(buf, pipeline.Script)
	if err != nil {
		panic(err)
	}
	script := buf.String()
	slog.Info("build script", "script", script, "pipeline name", pipeline.Name, "script type", pipeline.EventName)
	go func() {
		timer := time.NewTimer(s.Timeout)
		result := s.execute(script, args...)
		select {
		case <-timer.C:
			r <- &CodeResult{Status: ErrorStatus, Output: "timeout "}
		case <-result:
			r = result
		}
	}()
	return r
}

func (s *ScriptEngine) execute(script string, args ...interface{}) chan *CodeResult {
	result := make(chan *CodeResult)
	var wg conc.WaitGroup
	wg.Go(func() {
		codeResult, err := s.executeScript(script, args...)
		if err != nil {
			result <- &CodeResult{Status: ErrorStatus, Output: err}
		} else {
			result <- codeResult
		}
	})
	return result
}

func (s *ScriptEngine) executeScript(script string, args ...interface{}) (*CodeResult, error) {
	// build script engine
	input := make(map[string]any)
	vm := otto.New()
	_ = vm.Set("input", input)
	if args != nil {
		_ = vm.Set("args", args)
	}
	_ = vm.Set("send", func(call otto.FunctionCall) otto.Value {
		out := make(map[string]any)
		bytes, err := call.Argument(0).MarshalJSON()
		if err != nil {
			out["status"] = ErrorStatus
			out["output"] = err
			r, _ := vm.ToValue(out)
			return r
		}
		in := make(map[string]any)
		if err := json.Unmarshal(bytes, &in); err != nil {
			out["status"] = ErrorStatus
			out["output"] = err
			r, _ := vm.ToValue(out)
			return r
		}
		out["output"] = in
		status, _ := call.Argument(1).ToString()
		out["status"] = status
		r, _ := vm.ToValue(out)
		return r
	})
	// execute
	value, err := vm.Run(script)
	if err != nil {
		return nil, err
	}
	bytes, err := value.MarshalJSON()
	if err != nil {
		panic(err)
	}
	out := &CodeResult{}
	if err := json.Unmarshal(bytes, &out); err != nil {
		return nil, err
	}
	slog.Info("execute script success", "script", script, "output", out)
	return out, nil
}

func (s *ScriptEngine) DispatchById(pipelineId int64, args ...interface{}) <-chan *CodeResult {
	pipeline, err := s.PipelineService.GetPipelineById(pipelineId)
	if err != nil {
		slog.Error("dispatch pipeline script has error", "err", err)
		return nil
	}
	return s.Dispatch(pipeline, args...)
}

func (s *ScriptEngine) DispatchByEventKey(eventKey domain.EventKey, args ...interface{}) <-chan *CodeResult {
	pipelines := s.PipelineService.GetPipelineByEventKey(eventKey)
	for _, pipeline := range pipelines {
		s.Dispatch(pipeline, args...)
	}
	return nil
}

// DispatchLoginEvent dispatch login event
func (s *ScriptEngine) DispatchLoginEvent(args ...interface{}) <-chan *CodeResult {
	return s.DispatchByEventKey(domain.LoginEvent, args...)
}

// DispatchLogoutEvent dispatch logout event
func (s *ScriptEngine) DispatchLogoutEvent(args ...interface{}) <-chan *CodeResult {
	return s.DispatchByEventKey(domain.LogoutEvent, args...)
}

// DispatchBeforeUpdateArticleEvent dispatch before update article event
func (s *ScriptEngine) DispatchBeforeUpdateArticleEvent(args ...interface{}) <-chan *CodeResult {
	return s.DispatchByEventKey(domain.BeforeUpdateArticleEvent, args...)
}

// DispatchAfterUpdateArticleEvent dispatch after update article event
func (s *ScriptEngine) DispatchAfterUpdateArticleEvent(args ...interface{}) <-chan *CodeResult {
	return s.DispatchByEventKey(domain.AfterUpdateArticleEvent, args...)
}

// DispatchDeleteArticleEvent dispatch delete article event
func (s *ScriptEngine) DispatchDeleteArticleEvent(args ...interface{}) <-chan *CodeResult {
	return s.DispatchByEventKey(domain.DeleteArticleEvent, args...)
}

// DispatchBeforeUpdateDraftEvent dispatch before update draft event
func (s *ScriptEngine) DispatchBeforeUpdateDraftEvent(args ...interface{}) <-chan *CodeResult {
	return s.DispatchByEventKey(domain.BeforeUpdateDraftEvent, args...)
}

// DispatchAfterUpdateDraftEvent dispatch after update draft event
func (s *ScriptEngine) DispatchAfterUpdateDraftEvent(args ...interface{}) <-chan *CodeResult {
	return s.DispatchByEventKey(domain.AfterUpdateDraftEvent, args...)
}

// DispatchDeleteDraftEvent dispatch delete draft event
func (s *ScriptEngine) DispatchDeleteDraftEvent(args ...interface{}) <-chan *CodeResult {
	return s.DispatchByEventKey(domain.DeleteDraftEvent, args...)
}

// DispatchUpdateSiteInfoEvent dispatch update site info event
func (s *ScriptEngine) DispatchUpdateSiteInfoEvent(args ...interface{}) <-chan *CodeResult {
	return s.DispatchByEventKey(domain.UpdateSiteInfoEvent, args...)
}

// DispatchManualTriggerEvent dispatch manual trigger event
func (s *ScriptEngine) DispatchManualTriggerEvent(args ...interface{}) <-chan *CodeResult {
	return s.DispatchByEventKey(domain.ManualTriggerEvent, args...)
}
