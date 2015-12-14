package main

import "testing"

type StuntTransmitter struct {
	TransmitCalled bool
}

func (st *StuntTransmitter) Transmit() {
	st.TransmitCalled = true
}

func TestNotifier(t *testing.T) {
	st := StuntTransmitter{}
	n := NewNotifier([]Transmitter{&st})
	n.Start()
	alert := Alert{
		Name: "foo",
	}
	n.GetChan() <- &alert
	if !st.TransmitCalled {
		t.Error("transmit not called")
	}
}
