package main

import (
	"flag"
	"log"

	"server_demo/api/profile"
	"server_demo/internal/service"
	"server_demo/pkg/app"
)

func main() {
	optConfig := flag.String("f", "../../configs/profile/app_dev.yaml", "Config file path")
	flag.Parse()

	app, err := app.InitializeApp(*optConfig)
	if err != nil {
		log.Fatal(err)
	}

	service, err := service.InitializeService(app.ConfRaw)
	if err != nil {
		log.Fatal(err)
	}
	profile.RegisterServiceServer(app.GrpcServer, service)

	err = app.Start()
	if err != nil {
		log.Fatal(err)
	}
}
