package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type StuntListener struct {
	SensuChan chan *SensuResult
}

func (sl *StuntListener) Start() {
}

func (sl *StuntListener) GetSensuChan() chan *SensuResult {
	return sl.SensuChan
}

func TestAlertsHandler(t *testing.T) {
	recorder := httptest.NewRecorder()
	jw := &JsonResponseWriter{w: recorder}
	r, err := http.NewRequest("GET", "url", strings.NewReader("{}"))
	if err != nil {
		t.Error("error creating http.Request")
	}
	sl := StuntListener{SensuChan: make(chan *SensuResult, 10)}
	TestAlert(jw, r, &sl)

	if recorder.Code != 201 || recorder.Body.String() != "{\"ok\":1}" {
		t.Errorf("Expected %s, was %s", "201 - {\"ok\":1}", fmt.Sprintf("%d - %s", recorder.Code, recorder.Body.String()))
	}
}
