package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
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
	Webhook    string
	Secret     string
	HTTPClient *http.Client
}

func New(webhook, secret string) *DingTalkClient {
	return &DingTalkClient{
		Webhook: webhook,
		Secret:  secret,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
				DisableKeepAlives:   false,
			},
		},
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

	// 重试机制:最多重试3次
	var lastErr error
	for i := 0; i < 3; i++ {
		if i > 0 {
			// 等待一段时间后重试,使用指数退避
			fmt.Printf("第%d次重试,等待%dms...\n", i+1, i*500)
			time.Sleep(time.Duration(i*500) * time.Millisecond)
		}

		// 每次重试都重新生成签名和URL(避免时间戳过期)
		URL := c.Webhook
		if c.Secret != "" {
			timestamp, sign, _ := c.Sign()
			URL += fmt.Sprintf("&timestamp=%s&sign=%s", timestamp, sign)
			fmt.Printf("请求URL(第%d次): %s\n", i+1, URL)
		}

		req, err := http.NewRequest("POST", URL, bytes.NewBuffer(body))
		if err != nil {
			lastErr = fmt.Errorf("创建请求失败: %w", err)
			fmt.Printf("创建请求失败(第%d次): %v\n", i+1, err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Close = true // 禁用连接复用,避免陈旧连接

		fmt.Printf("开始发送请求(第%d次)...\n", i+1)
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("发送请求失败(第%d次): %w", i+1, err)
			fmt.Printf("发送请求失败(第%d次): %v\n", i+1, err)
			// EOF错误通常是临时性网络问题,继续重试
			continue
		}

		// 读取响应体
		respBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()

		if readErr != nil {
			lastErr = fmt.Errorf("读取响应失败(第%d次): %w", i+1, readErr)
			continue
		}

		// 打印响应信息
		fmt.Printf("钉钉API响应 [状态码: %d]:\n%s\n", resp.StatusCode, string(respBody))

		// 检查HTTP状态码
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("钉钉API返回错误状态码(第%d次): %d, 响应: %s", i+1, resp.StatusCode, string(respBody))
			continue
		}

		// 解析响应JSON,检查钉钉API返回的errcode
		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err == nil {
			if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
				errmsg := result["errmsg"]
				lastErr = fmt.Errorf("钉钉API返回错误(第%d次): errcode=%v, errmsg=%v", i+1, errcode, errmsg)
				continue
			}
		}

		// 成功
		return nil
	}

	return lastErr
}
