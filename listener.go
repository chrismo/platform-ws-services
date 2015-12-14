package main

type IListener interface {
	GetChan() chan *Alert
	Start()
}
