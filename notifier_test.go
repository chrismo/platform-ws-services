package main

import "testing"

type StuntTransmitter struct {
	TransmitCalled bool
	AlertPackage   AlertPackage
}

func (st *StuntTransmitter) Transmit(p AlertPackage) *TransmitResult {
	st.TransmitCalled = true
	p.Alert.Name = p.Alert.Name + "-done"
	st.AlertPackage = p
	return nil
}

func TestNotifierSimple(t *testing.T) {
	setupTestDB()
	defer tearDownTestDB()

	setupFakeDeployment()

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
	equals(t, "foobox0-check-name-done", st.AlertPackage.Alert.Name)
	equals(t, "1", st.AlertPackage.Deployment.Id)
	equals(t, "check-name", st.AlertPackage.Check.Name)
	equals(t, *new(Settings), st.AlertPackage.Settings)
}

func TestNotifierAlertPackageByValue(t *testing.T) {
	// multiple transmitters shouldn't be to alter data for each other
	setupTestDB()
	defer tearDownTestDB()

	setupFakeDeployment()

	stA, stB := StuntTransmitter{}, StuntTransmitter{}
	n := NewNotifier([]Transmitter{&stA, &stB})
	alert := Alert{
		Name:         "foobox0-check-name",
		DeploymentID: "1",
		CapsuleName:  "foobox0",
	}
	n.GetChan() <- &alert
	n.listenForAlert()

	equals(t, "foobox0-check-name-done", stA.AlertPackage.Alert.Name)
	equals(t, "foobox0-check-name", alert.Name)
}

func setupFakeDeployment() {
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
}
