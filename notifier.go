package main

import (
	"errors"
	"log"
	"strings"
	"time"
)

// Transmitter interface can be used to send an alert to any external service:
// SMS, PagerDuty, Slack, etc.
type Transmitter interface {
	Transmit()
}

// Notifier sends notifications for alerts based on deployment/group settings
type Notifier struct {
	NotifyChan   chan *Alert
	transmitters []Transmitter
}

// NewNotifier creates a new Notifier
func NewNotifier(transmitters []Transmitter) *Notifier {
	return &Notifier{
		NotifyChan:   make(chan *Alert),
		transmitters: transmitters,
	}
}

// Start starts monitoring
func (n *Notifier) Start() {
	go n.listenForAlerts()
}

// GetChan returns the channel to communication alerts
func (n *Notifier) GetChan() chan *Alert {
	return n.NotifyChan
}

func (n *Notifier) listenForAlerts() {
	for {
		select {
		case result := <-n.NotifyChan:
			if err := n.processAlert(result); err != nil {
				log.Printf("ERROR: Unable to process alert, %s", err.Error())
			}
		case <-time.After(100 * time.Millisecond):
			// NOP, just breath
		}
	}
}

func (n *Notifier) processAlert(alert *Alert) error {
	for _, transmitter := range n.transmitters {
		// TODO: each of these as goroutine?
		transmitter.Transmit()
	}
	return nil
}

func getInfo(alert *Alert) (Deployment, Check, Settings, error) {
	latterPart := strings.TrimPrefix(alert.Name, alert.CapsuleName)
	checkName := strings.Replace(latterPart, "-", "", 1)
	var deployment Deployment
	var check Check
	var settings Settings
	deployment, err := LookupDeploymentById(alert.DeploymentID)
	if err != nil {
		// exceptions, shmexceptions
		return deployment, check, settings, err
	} else if check, err := deployment.CheckByName(checkName); err == nil {
		settings, err := deployment.MergedAlertSettings()
		if err != nil {
			return deployment, check, settings, err
		}
		return deployment, check, settings, nil
	} else {
		return deployment, check, settings, errors.New("check not registered for alerts")
	}
}
