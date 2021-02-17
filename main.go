package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/websublime/courier/api"
	"github.com/websublime/courier/config"
	"github.com/websublime/courier/storage"
)

func main() {
	boot()
}

func boot() {
	env := config.LoadEnvironmentConfig()

	db, err := storage.Dial(env)
	if err != nil {
		logrus.Fatalf("Error opening database: %+v", err)
	}
	defer db.Close()

	app := config.BootApplication()

	api.WithVersion(app, env, db)

	app.Listen(fmt.Sprintf("%s:%s", env.CourierHost, env.CourierPort))
}
