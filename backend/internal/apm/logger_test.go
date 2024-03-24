package apm

import (
	"github.com/asaskevich/EventBus"
	"github.com/stretchr/testify/assert"
	"testing"
)

var config = Config{Url: "http://43.143.195.208:5080", Authorization: "Basic amlhbmd3MTAyN0BnbWFpbC5jb206YWxsaW9AMjAyMw==", Organization: "default", Stream: "fusion"}
var l = NewLogger(EventBus.New(), &config)

func TestPushToOpenObserve(t *testing.T) {
	Init(l)
	l.Put("[{\"level\":\"info\",\"job\":\"test\",\"log\":\"test message for openobserve\"}]")

	l.bus.WaitAsync()
}

func TestQueryOpenObserve(t *testing.T) {

	a := assert.New(t)
	Init(l)
	search := NewSearch()
	search.SQL("select * from fusion")
	f, err := l.FindMap(search)
	a.NoError(err)

	a.NotNil(f)

	result, err := f.Format()
	a.NoError(err)

	a.NotNil(result)
}
