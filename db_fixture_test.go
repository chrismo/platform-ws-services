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
	if b, _ := dbExists(name); !b {
		r.DBCreate(name).RunWrite(session)
	}

	for _, tableName := range tableNames() {
		if b, _ := tableExists(name, tableName); !b {
			r.DB(name).TableCreate(tableName).RunWrite(session)
		}
	}

	if b, _ := indexExists(name, "checks", "type_name"); !b {
		_, err := r.DB(name).Table("checks").IndexCreateFunc("type_name", func(row r.Term) interface{} {
			return []interface{}{row.Field("type"), row.Field("name")}
		}).RunWrite(session)
		if err != nil {
			log.Fatalf("Error creating index: %s", err)
		}
	}

	if b, _ := indexExists(name, "checks", "type"); !b {
		_, err := r.DB(name).Table("checks").IndexCreate("type").RunWrite(session)
		if err != nil {
			log.Fatalf("Error creating index: %s", err)
		}
	}
	session.Use(name)
}

func tearDownDB(name string, session *r.Session) {
	for _, tableName := range tableNames() {
		// just truncating the tables keeps tests fast
		r.DB(name).Table(tableName).Delete().RunWrite(session)
	}
}

func tableNames() []string {
	return []string{"deployments", "checks", "groups"}
}

func dbExists(name string) (bool, error) {
	cur, _ := r.DBList().Run(session)
	defer cur.Close()
	return existsInListFromDB(name, cur)
}

func tableExists(dbName, tableName string) (bool, error) {
	cur, _ := r.DB(dbName).TableList().Run(session)
	defer cur.Close()
	return existsInListFromDB(tableName, cur)
}

func indexExists(dbName, tableName, indexName string) (bool, error) {
	cur, _ := r.DB(dbName).Table(tableName).IndexList().Run(session)
	defer cur.Close()
	return existsInListFromDB(indexName, cur)
}

func existsInListFromDB(lookFor string, cursor *r.Cursor) (bool, error) {
	names := []string{}
	cursor.All(&names)
	return contains(names, lookFor), nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
