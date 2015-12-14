package main

import (
	"errors"
	"log"
	"strings"
	"time"
)

// AlertPackage has all related structs bundled up so a Transmitter can decide
// for itself what to transmit.
type AlertPackage struct {
	Alert
	Deployment
	Check
	Settings
}

// Transmitter interface can be used to send an alert to any external service:
// SMS, PagerDuty, Slack, etc.
type Transmitter interface {
	Transmit(p AlertPackage)
}

// Notifier sends notifications for alerts based on deployment/group settings
type Notifier struct {
	NotifyChan   chan *Alert
	transmitters []Transmitter
}

// NewNotifier creates a new Notifier
func NewNotifier(transmitters []Transmitter) *Notifier {
	return &Notifier{
		NotifyChan:   make(chan *Alert, 10),
		transmitters: transmitters,
	}
}

// Start starts monitoring
func (n *Notifier) Start() {
	go func() {
		for {
			n.listenForAlert()
		}
	}()
}

// GetChan returns the channel to communication alerts
func (n *Notifier) GetChan() chan *Alert {
	return n.NotifyChan
}

func (n *Notifier) listenForAlert() {
	select {
	case result := <-n.NotifyChan:
		if err := n.processAlert(result); err != nil {
			log.Printf("ERROR: Unable to process alert, %s", err.Error())
		}
	case <-time.After(100 * time.Millisecond):
		// NOP, just breath
	}
}

func (n *Notifier) processAlert(alert *Alert) error {
	alertPackage, err := buildAlertPackage(alert)
	if err != nil {
		return err
	}
	for _, transmitter := range n.transmitters {
		// TODO: each of these as goroutine so one bad one can't block the rest?
		transmitter.Transmit(*alertPackage)
	}
	return nil
}

func buildAlertPackage(alert *Alert) (*AlertPackage, error) {
	latterPart := strings.TrimPrefix(alert.Name, alert.CapsuleName)
	checkName := strings.Replace(latterPart, "-", "", 1)
	deployment, err := LookupDeploymentById(alert.DeploymentID)
	if err != nil {
		return nil, err
	} else if check, err := deployment.CheckByName(checkName); err == nil {
		settings, err := deployment.MergedAlertSettings()
		if err != nil {
			return nil, err
		}
		return &AlertPackage{
			Alert:      *alert,
			Deployment: deployment,
			Check:      check,
			Settings:   settings,
		}, nil
	} else {
		return nil, errors.New("check not registered for alerts")
	}
}
