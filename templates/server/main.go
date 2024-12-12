package main

import (
	server "github.com/RedeployAB/go-template/templates/server/server"
)

func main() {
	log := server.NewLogger()

	srv := server.New(server.WithOptions(server.Options{
		Logger: log,
	}))

	if err := srv.Start(); err != nil {
		log.Error("Server error.", "error", err)
	}
}
