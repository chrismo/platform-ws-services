package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	pagerdutyURL = "https://events.pagerduty.com/generic/2010-04-15/create_event.json"
)

// Pagerduty model implements Transmitter for handling PagerDuty events.
type Pagerduty struct {
}

// PagerdutyEvent implements the payload for the PD events API.
type PagerdutyEvent struct {
	ServiceKey  string `json:"service_key"`            //	The GUID of one of your "Generic API" services. This is the "service key" listed on a Generic API's service detail page.
	EventType   string `json:"event_type"`             //	Set this to "trigger" or "resolve"
	IncidentKey string `json:"incident_key,omitempty"` // Required for resolve. If not used in trigger, will be auto-assigned one.
	Description string `json:"description,omitempty"`  // Required for trigger. Max 1024 char description
	// TODO: include optional Details JSON object/struct
}

// Transmit handles an AlertPackage and creates/resolves PagerDuty incidents from it.
func (pd *Pagerduty) Transmit(ap AlertPackage) *TransmitResult {
	var msg string

	if len(ap.Settings.PagerdutyKey) == 0 {
		return &TransmitResult{Result: Skipped, Message: "Deployment/Group has no PagerDuty Setting configured."}
	}

	e := &PagerdutyEvent{
		ServiceKey:  ap.Settings.PagerdutyKey,
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
		// TODO: Need to ID the alert
		msg = "ERROR: Unable to marshal payload to transmit to PagerDuty"
		log.Print(msg)
		return &TransmitResult{Result: Error, Message: msg}
	}

	resp, err := http.Post(pagerdutyURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		// TODO: Need to ID the alert
		msg = "ERROR: Unable to transmit alert to PagerDuty"
		log.Print(msg)
		return &TransmitResult{Result: Error, Message: msg}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// TODO: Need to ID the alert
		var body string
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			body = fmt.Sprintf("[Unable to read response body: %s]", err.Error())
		} else {
			body = string(bodyBytes)
		}

		msg = fmt.Sprintf("ERROR: Post to PagerDuty unsuccessful: %d - %s", resp.StatusCode, body)
		log.Print(msg)
		return &TransmitResult{Result: Error, Message: msg}
	}
	return &TransmitResult{Result: Success, Message: fmt.Sprintf("PagerDuty %s event posted.", e.EventType)}

}
