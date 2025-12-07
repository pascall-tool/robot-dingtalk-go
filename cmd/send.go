package cmd


import (
	"strings"

	"github.com/spf13/cobra"
	"dingtalk-action/pkg/dingtalk"
)

var (
	webhook   string
	secret    string
	message   string
	atMobiles string
	markdown  bool
	title     string
)
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "发送钉钉消息",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := dingtalk.New(webhook, secret)

		mobiles := []string{}
		if atMobiles != "" {
			mobiles = strings.Split(atMobiles, ",")
		}

		var err error
		if markdown {
			if title == "" {
				title = "Notification"
			}
			err = client.SendMarkdown(title, message, mobiles)
		} else {
			err = client.SendText(message, mobiles)
		}

		if err == nil {
			// 成功日志
			println("发送钉钉消息成功！")
			println("dingtalk-action version:", Version)
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().StringVar(&webhook, "webhook", "", "钉钉 Webhook 地址")
	sendCmd.Flags().StringVar(&secret, "secret", "", "加签密钥")
	sendCmd.Flags().StringVar(&message, "msg", "", "消息内容")
	sendCmd.Flags().StringVar(&atMobiles, "at", "", "@手机号列表（逗号分隔）")
	sendCmd.Flags().BoolVar(&markdown, "md", false, "使用 Markdown 消息")
	sendCmd.Flags().StringVar(&title, "title", "Notification", "Markdown 标题")

	sendCmd.MarkFlagRequired("webhook")
	sendCmd.MarkFlagRequired("msg")
}
