package main

import (
	"github.com/RedeployAB/go-template/templates/service/service"
)

func main() {
	log := service.NewLogger()

	svc := service.New(service.WithOptions(service.Options{
		Logger: log,
	}))

	if err := svc.Start(); err != nil {
		log.Error("Service error.", "error", err)
	}
}
