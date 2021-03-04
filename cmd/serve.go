package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/websublime/courier/api"
	"github.com/websublime/courier/config"
	"github.com/websublime/courier/storage"
)

var serveCmd = cobra.Command{
	Use:  "serve",
	Long: "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, serve)
	},
}

func serve(conf *config.EnvironmentConfig) {
	db, err := storage.Dial(conf)
	if err != nil {
		logrus.Fatalf("Error opening database: %+v", err)
	}
	defer db.Close()

	app := config.BootApplication(conf)

	api.WithVersion(app, conf, db)

	app.Listen(fmt.Sprintf("%s:%s", conf.CourierHost, conf.CourierPort))
}
