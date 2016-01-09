package main

import (
	"log"

	r "github.com/dancannon/gorethink"
)

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

	_, err := r.DB(name).Table("checks").IndexCreateFunc("type_name", func(row r.Term) interface{} {
		return []interface{}{row.Field("type"), row.Field("name")}
	}).RunWrite(session)
	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}

	_, err = r.DB(name).Table("checks").IndexCreate("type").RunWrite(session)
	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}

	session.Use(name)
}

func tearDownDB(name string, session *r.Session) {
	r.DBDrop(name).RunWrite(session)
}
