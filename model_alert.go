package main

import (
	"encoding/json"
	"io"
)

const (
	Resolved = iota // 0
	Warning         // 1
	Critical        // 2
	Unknown         // 3
)

type ClientAlert struct {
	Client string `json:"client"` // the client the check came from
	Alert  Alert  `json:"check"`  // the check info
}

// Alert is an instance of a failed check of a deployment type.
type Alert struct {
	Name         string `json:"name"`
	CapsuleName  string `json:"capsule_name"`
	Output       string `json:"output"`
	Status       int    `json:"status"`
	CapsuleID    string `json:"capsule_id,omitempty"`
	DeploymentID string `json:"deployment_id,omitempty"`
	AccountSlug  string `json:"account,omitempty"`
}

// NewAlertFromJSON generates a new alert from a json string
func NewAlertFromJSON(r io.Reader) (*Alert, error) {
	decoder := json.NewDecoder(r)
	var clientAlert ClientAlert
	err := decoder.Decode(&clientAlert)
	if err != nil {
		return nil, err
	}
	return &clientAlert.Alert, nil
}

func (a *Alert) Serialize() ([]byte, error) {
	return json.Marshal(a)
}
