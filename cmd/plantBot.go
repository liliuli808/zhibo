package cmd

import (
	"github.com/spf13/cobra"
	"time"
	"zhibo/server"
)

func init() {
	var serverCmd = &cobra.Command{
		Use:   "plantBot",
		Short: "plantBot",
		Run: func(cmd *cobra.Command, args []string) {
			sendPlantBot()
		},
	}
	serverCmd.PersistentFlags().StringVar(&botConfig, "config", "./plantBot.yaml", "配置文件地址")
	rootCmd.AddCommand(serverCmd)
}

func sendPlantBot() {
	serverConfig := server.GetConfig(config)
	agent := server.NewAgent(serverConfig)
	for {
		agent.StartSendPlantBot()
		time.Sleep(time.Second * 5)
	}
}
