package main

import (
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
	alert, err := NewAlertFromJSON(r.Body)
	if err != nil {
		w.WriteError(err)
		return
	}
	listener.GetSensuChan() <- alert
	w.WriteOk(201)
}
