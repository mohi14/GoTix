package main

import (
	"gotickets/internal/config"
	"gotickets/internal/server"
)

func main() {
	// load environment variables
	cfg := config.LoadEnv()

	// connect to the database
	db := config.ConnectDB(cfg)

	server.StartServer(db, cfg)
}
