package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type StuntListener struct {
	SensuChan chan *Alert
}

func (sl *StuntListener) Start() {
}

func (sl *StuntListener) GetSensuChan() chan *Alert {
	return sl.SensuChan
}

func TestAlertsHandler(t *testing.T) {
	// TODO: this is only testing an empty body. Other tests needed.
	recorder := httptest.NewRecorder()
	jw := &JsonResponseWriter{w: recorder}
	r, err := http.NewRequest("GET", "url", strings.NewReader("{}"))
	if err != nil {
		t.Error("error creating http.Request")
	}
	sl := StuntListener{SensuChan: make(chan *Alert, 10)}
	TestAlert(jw, r, &sl)

	if recorder.Code != 201 || recorder.Body.String() != "{\"ok\":1}" {
		t.Errorf("Expected %s, was %s", "201 - {\"ok\":1}", fmt.Sprintf("%d - %s", recorder.Code, recorder.Body.String()))
	}
}
