package main

import (
	"encoding/json"
	"io"
)

// Alert is an instance of a failed check of a deployment type.
type Alert struct {
	Name         string  `json:"name"`
	CapsuleName  string  `json:"capsule_name"`
	Output       string  `json:"output"`
	Status       float64 `json:"status"`
	CapsuleID    string  `json:"capsule_id,omitempty"`
	DeploymentID string  `json:"deployment_id,omitempty"`
	AccountSlug  string  `json:"account,omitempty"`
}

// NewAlertFromJSON generates a new alert from a json string
func NewAlertFromJSON(r io.Reader) (*Alert, error) {
	decoder := json.NewDecoder(r)
	var alert Alert
	err := decoder.Decode(&alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

func (a *Alert) serialize() ([]byte, error) {
	return json.Marshal(a)
}
