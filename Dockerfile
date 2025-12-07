FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY . .
RUN go build -o dingtalk-action .

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/dingtalk-action /usr/local/bin/dingtalk-action

ENTRYPOINT ["dingtalk-action"]
