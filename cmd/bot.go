package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"time"
	"zhibo/bot"
)

var botConfig string

func init() {
	var serverCmd = &cobra.Command{
		Use:   "bot",
		Short: "bot",
		Run: func(cmd *cobra.Command, args []string) {
			for {
				send()
				time.Sleep(time.Second * 6)
			}
		},
	}
	serverCmd.PersistentFlags().StringVar(&botConfig, "config", "./bot.yaml", "配置文件地址")
	rootCmd.AddCommand(serverCmd)
}

func send() {
	bitConfig, err := bot.NewConfigWithFile(botConfig)
	if err != nil {
		panic(err)
	}
	bot.Instance(bitConfig).Consumer.Start(context.Background())
}
