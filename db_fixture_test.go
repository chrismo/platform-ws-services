package main

import r "github.com/dancannon/gorethink"

func setupTestDB() {
	session, _ = initRethinkConn()
	setupDB("alerts_test", session)
}

func tearDownTestDB() {
	tearDownDB("alerts_test", session)
}

func setupDB(name string, session *r.Session) {
	r.DBCreate(name).RunWrite(session)
	r.DB(name).TableCreate("deployments").RunWrite(session)
	r.DB(name).TableCreate("checks").RunWrite(session)
	r.DB(name).TableCreate("groups").RunWrite(session)
	session.Use(name)
}

func tearDownDB(name string, session *r.Session) {
	r.DBDrop(name).RunWrite(session)
}
