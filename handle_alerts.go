package main

import (
	"net/http"
	"strings"
)

// AlertsHandlerFunc creates an anonymous function to handle alert POSTs
func AlertsHandlerFunc(listeners []Listener) HandlerFunc {
	return func(w *JsonResponseWriter, r *http.Request) {
		if strings.ToUpper(r.Method) == "POST" {
			TestAlert(w, r, listeners)
		} else {
			UnknownEndpoint(w, r)
		}
	}
}

// TestAlert processes a POSTed alert
func TestAlert(w *JsonResponseWriter, r *http.Request, listeners []Listener) {
	alert, err := NewAlertFromJSON(r.Body)
	if err != nil {
		w.WriteError(err)
		return
	}
	for _, l := range listeners {
		l.GetChan() <- alert
	}

	w.WriteOk(201)
}
