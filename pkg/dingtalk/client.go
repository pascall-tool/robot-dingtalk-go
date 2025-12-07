package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type MessageType string

const (
	Text     MessageType = "text"
	Markdown MessageType = "markdown"
)

type DingTalkClient struct {
	Webhook string
	Secret  string
}

func New(webhook, secret string) *DingTalkClient {
	return &DingTalkClient{
		Webhook: webhook,
		Secret:  secret,
	}
}

// 钉钉签名（如果 secret 存在）
func (c *DingTalkClient) Sign() (timestamp string, sign string, err error) {
	if c.Secret == "" {
		return "", "", nil
	}

	timestamp = fmt.Sprintf("%d", time.Now().UnixMilli())
	stringToSign := timestamp + "\n" + c.Secret

	h := hmac.New(sha256.New, []byte(c.Secret))
	h.Write([]byte(stringToSign))
	signBytes := h.Sum(nil)

	sign = base64.StdEncoding.EncodeToString(signBytes)
	return timestamp, url.QueryEscape(sign), nil
}

func (c *DingTalkClient) SendText(message string, atMobiles []string) error {
	payload := map[string]any{
		"msgtype": "text",
		"text": map[string]string{
			"content": message,
		},
		"at": map[string]any{
			"atMobiles": atMobiles,
			"isAtAll":   false,
		},
	}

	return c.send(payload)
}

func (c *DingTalkClient) SendMarkdown(title, text string, atMobiles []string) error {
	payload := map[string]any{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  text,
		},
		"at": map[string]any{
			"atMobiles": atMobiles,
			"isAtAll":   false,
		},
	}

	return c.send(payload)
}

func (c *DingTalkClient) send(payload any) error {
	body, _ := json.Marshal(payload)

	URL := c.Webhook

	// 带签名
	if c.Secret != "" {
		timestamp, sign, _ := c.Sign()
		URL += fmt.Sprintf("&timestamp=%s&sign=%s", timestamp, sign)
	}

	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
