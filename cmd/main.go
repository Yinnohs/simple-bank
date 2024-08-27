package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/yinnohs/simple-bank/api"
	db "github.com/yinnohs/simple-bank/db/sqlc"
)

const (
	dbDriver    = "postgres"
	dbSource    = "postgresql://yinnohs:1234@localhost:5432/simple_bank?sslmode=disable"
	baseAddress = "127.0.0.1:5050"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err.Error())
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(baseAddress)
	if err != nil {
		log.Fatal("Cannot start the server for the next REASON: ", err.Error())
	}
}
