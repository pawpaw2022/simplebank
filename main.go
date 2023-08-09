package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // postgresql driver
	"github.com/pawpaw2022/simplebank/api"
	db "github.com/pawpaw2022/simplebank/db/postgresql"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddr = "localhost:8080"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: %w", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddr)

	if err != nil {
		log.Fatal("cannot start server: %w", err)
	}

}
