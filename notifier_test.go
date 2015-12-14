package main

import "testing"

type StuntTransmitter struct {
	TransmitCalled bool
	AlertPackage   AlertPackage
}

func (st *StuntTransmitter) Transmit(p AlertPackage) {
	st.TransmitCalled = true
	st.AlertPackage = p
}

func TestNotifier(t *testing.T) {
	setupTestDB()
	defer tearDownTestDB()

	d := Deployment{
		Id:      "1",
		GroupId: "1",
		Type:    "foobox",
		Name:    "deployment-foobox",
		Checks: []Check{Check{
			Id:          "1",
			Type:        "check-type",
			Name:        "check-name",
			Level:       2,
			Title:       "Generic Check Title",
			Description: "Generic Check Description",
		}},
	}
	d.Save()

	st := StuntTransmitter{}
	n := NewNotifier([]Transmitter{&st})
	alert := Alert{
		Name:         "foobox0-check-name",
		DeploymentID: "1",
		CapsuleName:  "foobox0",
	}
	n.GetChan() <- &alert
	n.listenForAlert()

	assert(t, st.TransmitCalled, "Transmit in mock not called")
	equals(t, "foobox0-check-name", st.AlertPackage.Alert.Name)
	equals(t, "1", st.AlertPackage.Deployment.Id)
	equals(t, "check-name", st.AlertPackage.Check.Name)
	equals(t, *new(Settings), st.AlertPackage.Settings)
}
