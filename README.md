# dingtalk-action

钉钉机器人消息发送 CLI 工具，支持本地命令行和 GitHub/Gitea Actions 场景，适合团队自动化通知。

## 🚀 快速开始

### 1. 本地编译 & 运行

```bash
# 克隆仓库
# git clone <repo-url>
cd dingtalk-action

# 下载依赖（如在中国大陆，建议配置代理）
go mod tidy

# 编译
go build -o dingtalk-action .

# 发送钉钉消息（text）
./dingtalk-action send \
  --webhook <你的Webhook> \
  --secret <你的Webhook加签> \
  --msg "构建成功" \
  --at "13800000000,13900000000"

# 发送 Markdown 消息
./dingtalk-action send \
  --webhook <你的Webhook> \
  --secret <你的Webhook加签> \
  --md \
  --title "通知" \
  --msg "### 🚀 构建成功\n- 项目：xxx\n- 时间：$(date)"
```

### 2. Docker 构建 & 使用

```bash
# 构建镜像
# 推荐使用代理加速依赖下载
# docker build -t dingtalk-action .

# 运行
# docker run --rm dingtalk-action send --webhook <你的Webhook>   --secret <你的Webhook加签>  --msg "Hello"
```

### 3. GitHub/Gitea Actions 集成

在 workflow 中添加：

```yaml
- name: Send DingTalk Notification
  uses: mingcai-toolkit/dingtalk-action@v1
  with:
    webhook: ${{ secrets.DINGTALK_WEBHOOK }}
    secret:  ${{ secrets.DINGTALK_SECRET }}
    message: "构建成功 🎉"
    at_mobiles: "13800000000"
```

## ⚙️ 参数说明

| 参数         | 说明                 | 必填 | 示例                                                     |
|--------------|----------------------|------|--------------------------------------------------------|
| --webhook    | 钉钉机器人Webhook    | 是   | https://oapi.dingtalk.com/robot/send?access_token=xxxx |
| --secret     | 加签密钥             | 否   | xxxxx                                                  |
| --msg        | 消息内容             | 是   | "构建成功"                                                 |
| --at         | @手机号（逗号分隔）  | 否   | "13800000000,13900000000"                              |
| --md         | 使用Markdown消息     | 否   | true/false                                             |
| --title      | Markdown标题         | 否   | "通知"                                                   |

## 📝 进阶
- 支持 Docker 镜像发布到 ghcr/Docker Hub
- 支持多平台编译（goreleaser 可选）
- 支持自定义消息模板、JSON输入等扩展

## 🛠️ 常见问题
- **依赖下载慢/失败**：建议配置 Go 代理和 SOCKS5 代理（如 hysteria）
- **Webhook/密钥安全**：建议通过环境变量或 CI/CD secrets 管理

---

如需更多功能或定制，欢迎提 issue！
