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
			for {
				agent()
				time.Sleep(time.Second * 5)
			}
		},
	}
	serverCmd.PersistentFlags().StringVar(&config, "config", "./config.yaml", "配置文件地址")
	rootCmd.AddCommand(serverCmd)
}

func agent() {
	serverConfig := server.GetConfig(config)
	agent := server.NewAgent(serverConfig)
	agent.Wg.Add(1)
	defer agent.Wg.Done()
	agent.Start()
	agent.Wg.Wait()
}
