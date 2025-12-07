package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dingtalk-action",
	Short: "钉钉机器人消息发送工具",
}

func Execute() error {
	return rootCmd.Execute()
}
