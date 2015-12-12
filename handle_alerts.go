package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func AlertsHandlerFunc(listener IListener) HandlerFunc {
	return func(w *JsonResponseWriter, r *http.Request) {
		if strings.ToUpper(r.Method) == "POST" {
			TestAlert(w, r, listener)
		} else {
			UnknownEndpoint(w, r)
		}
	}
}

func TestAlert(w *JsonResponseWriter, r *http.Request, listener IListener) {
	decoder := json.NewDecoder(r.Body)
	var sensuResult SensuResult
	err := decoder.Decode(&sensuResult)
	if err != nil {
		w.WriteError(err)
		return
	}
	listener.GetSensuChan() <- &sensuResult
	w.WriteOk(201)
}
