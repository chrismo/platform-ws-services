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
	slackURL = "https://slack.com/api/chat.postMessage"
)

type Slack struct {
	ApiKey  string `gorethink:"api_key" json:"api_key"`
	Channel string `gorethink:"channel" json:"channel"`
}

// SlackPostMessage implements the params for the postMessage API
type SlackPostMessage struct {
	Token   string `json:"token"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	AsUser  bool   `json:"as_user"`
}

// TODO: take the alert and make a nice Slack notification with some sort of
// colored status
func (s *Slack) Transmit(ap AlertPackage) *TransmitResult {
	var msg string
	var alertType string

	if len(ap.Settings.Slack.ApiKey) == 0 {
		return &TransmitResult{Result: Skipped, Message: "Deployment/Group has no Slack Setting configured."}
	}

	if ap.Alert.Status == Resolved {
		alertType = "Resolved"
	} else {
		alertType = "Triggered"
	}

	// TODO: best unique ID for the Alert? Capsule ID? That's never persisted, right?
	text := fmt.Sprintf("Alert for %s, Capsule ID %s - %s", ap.Check.Description, ap.CapsuleID, alertType)

	e := &SlackPostMessage{
		Token:   ap.Settings.Slack.ApiKey,
		Channel: ap.Settings.Slack.Channel,
		Text:    text,
		AsUser:  true,
	}

	payload, err := json.Marshal(e)
	if err != nil {
		// TODO: Need to ID the alert
		msg = "ERROR: Unable to marshal payload to transmit to Slack"
		log.Print(msg)
		return &TransmitResult{Result: Error, Message: msg}
	}

	resp, err := http.Post(slackURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		// TODO: Need to ID the alert
		msg = "ERROR: Unable to transmit alert to Slack"
		log.Print(msg)
		return &TransmitResult{Result: Error, Message: msg}
	}
	defer resp.Body.Close()

	var body string
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		body = fmt.Sprintf("[Unable to read response body: %s]", err.Error())
	} else {
		body = string(bodyBytes)
	}
	log.Print(body)

	if resp.StatusCode != 200 {
		// TODO: Need to ID the alert
		msg = fmt.Sprintf("ERROR: Post to Slack unsuccessful: %d - %s", resp.StatusCode, body)
		log.Print(msg)
		return &TransmitResult{Result: Error, Message: msg}
	}
	return &TransmitResult{Result: Success, Message: fmt.Sprintf("%s alert posted to Slack.", alertType)}
}

//func (s *Slack) SendToSlack(deploymentName, message string, status float64) {

//
