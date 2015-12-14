package main

import (
	"errors"
	"strings"
)

type Notifier struct {
}

func NewNotifier() *Notifier {
	return &Notifier{}
}

func (n *Notifier) Start() {

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
