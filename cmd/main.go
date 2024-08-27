package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/yinnohs/simple-bank/api"
	db "github.com/yinnohs/simple-bank/db/sqlc"
	"github.com/yinnohs/simple-bank/util"
)

// const (
// 	dbDriver    = "postgres"
// 	dbSource    = "postgresql://yinnohs:1234@localhost:5432/simple_bank?sslmode=disable"
// 	baseAddress = "127.0.0.1:5050"
// )

func main() {
	config, err := util.LoadConfig("./")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err.Error())
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.BaseAddress)
	if err != nil {
		log.Fatal("Cannot start the server for the next REASON: ", err.Error())
	}
}
