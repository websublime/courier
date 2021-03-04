package cmd

import (
	"github.com/spf13/cobra"

	"github.com/websublime/courier/config"
)

var configFile = ""

var rootCmd = cobra.Command{
	Use: "courier",
	Run: func(cmd *cobra.Command, args []string) {
		execWithConfig(cmd, serve)
	},
}

// RootCommand will setup and return the root command
func RootCommand() *cobra.Command {
	rootCmd.AddCommand(&serveCmd, &versionCmd)

	return &rootCmd
}

func execWithConfig(cmd *cobra.Command, fn func(conf *config.EnvironmentConfig)) {
	env := config.LoadEnvironmentConfig()

	fn(env)
}
