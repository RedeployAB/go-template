package main

import (
	"github.com/RedeployAB/go-template/templates/http-server/server"
)

func main() {
	log := server.NewDefaultLogger()

	srv := server.New(server.WithOptions(server.Options{
		Logger: log,
	}))

	if err := srv.Start(); err != nil {
		log.Error("Server error.", "error", err)
	}
}
