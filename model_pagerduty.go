package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	pagerdutyURL = "https://events.pagerduty.com/generic/2010-04-15/create_event.json"
	// TODO: This is a setting, d00f
	pagerdutyServiceKey = "545bf39a778d45b1b4160a7fd782fae9" // free, trial account
)

// Pagerduty model implements Transmitter for handling PagerDuty events.
type Pagerduty struct {
}

// PagerdutyEvent implements the payload for the PD events API.
type PagerdutyEvent struct {
	ServiceKey  string `json:"service_key"`            //	The GUID of one of your "Generic API" services. This is the "service key" listed on a Generic API's service detail page.
	EventType   string `json:"event_type"`             //	Set this to "trigger" or "resolve"
	IncidentKey string `json:"incident_key,omitempty"` // Required for resolve.
	Description string `json:"description,omitempty"`  // Required for trigger. Max 1024 char description
	// TODO: include optional Details JSON object/struct
}

// Transmit handles an AlertPackage and creates/resolves PagerDuty incidents from it.
func (pd *Pagerduty) Transmit(ap AlertPackage) {
	e := &PagerdutyEvent{
		ServiceKey:  pagerdutyServiceKey,
		IncidentKey: ap.Alert.CapsuleID,
		// TODO: best unique ID for the Alert? Capsule ID? That's never persisted, right?
		Description: fmt.Sprintf("Alert for %s, Capsule ID %s", ap.Check.Description, ap.CapsuleID),
	}

	if ap.Alert.Status == Resolved {
		e.EventType = "resolve"
	} else {
		e.EventType = "trigger"
	}

	payload, err := json.Marshal(e)
	if err != nil {
		log.Print("ERROR: Unable to marshal payload to transmit to PagerDuty")
		return
	}

	resp, err := http.Post(pagerdutyURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		// TODO: Best ID?
		log.Print("ERROR: Unable to transmit alert to PagerDuty")
		return
	}
	defer resp.Body.Close()
}
