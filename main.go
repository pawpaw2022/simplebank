package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // postgresql driver
	"github.com/pawpaw2022/simplebank/api"
	db "github.com/pawpaw2022/simplebank/db/postgresql"
	"github.com/pawpaw2022/simplebank/util"
)

func main() {
	// Load config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: %w", err)
	}

	// Connect to db
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: %w", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: %w", err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server: %w", err)
	}

}
