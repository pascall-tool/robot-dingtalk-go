package cmd

import (
	"github.com/spf13/cobra"
)


var Version = "v1.0.1" // 默认值，实际编译时可覆盖

var rootCmd = &cobra.Command{
	Use:   "dingtalk-action",
	Short: "钉钉机器人消息发送工具",
}

func Execute() error {
	return rootCmd.Execute()
}
