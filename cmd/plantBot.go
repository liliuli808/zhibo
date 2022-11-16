package cmd

import (
	"github.com/spf13/cobra"
	"time"
	"zhibo/server"
)

var plantConfig string

func init() {
	var serverCmd = &cobra.Command{
		Use:   "plantBot",
		Short: "plantBot",
		Run: func(cmd *cobra.Command, args []string) {
			sendPlantBot()
		},
	}
	serverCmd.PersistentFlags().StringVar(&plantConfig, "config", "./config.yaml", "配置文件地址")
	rootCmd.AddCommand(serverCmd)
}

func sendPlantBot() {
	serverConfig := server.GetConfig(plantConfig)
	agent := server.NewAgent(serverConfig)
	for {
		agent.StartSendPlantBot()
		time.Sleep(time.Second * 5)
	}
}
