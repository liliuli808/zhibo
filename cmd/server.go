package cmd

import (
	"github.com/spf13/cobra"
	"time"
	"zhibo/server"
)

var config string

func init() {
	var serverCmd = &cobra.Command{
		Use:   "server",
		Short: "server",
		Run: func(cmd *cobra.Command, args []string) {
			agent()
		},
	}
	serverCmd.PersistentFlags().StringVar(&config, "config", "./server.yaml", "配置文件地址")
	rootCmd.AddCommand(serverCmd)
}

func agent() {
	serverConfig := server.GetConfig(config)
	agent := server.NewAgent(serverConfig)
	for {
		agent.Start()
		time.Sleep(time.Second * 5)
	}
}
