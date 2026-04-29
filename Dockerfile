# 构建阶段：直接用 1ms.run 的 golang:1.24-alpine3.20
FROM docker.1ms.run/library/golang:1.24-alpine3.20 AS builder

# 保留 CGO + 静态编译
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOPROXY=https://goproxy.cn,direct

# 国内 APK 源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装编译依赖
RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .

# 静态编译 + 瘦身
RUN go build -ldflags="-w -s" -o school-app cmd/server/main.go

# 构建数据库初始化工具
RUN go build -ldflags="-w -s" -o init-db scripts/init-db.go

# 运行阶段：直接用 1ms.run 的 alpine:3.20
FROM docker.1ms.run/library/alpine:3.20

# 国内 APK 源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装必要证书（最小化）
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/school-app /app/
COPY --from=builder /app/init-db /app/
COPY static/ /app/static/
COPY templates/ /app/templates/

# 复制数据库文件（如果存在）
COPY school.db /app/school.db 2>/dev/null || true

# 复制环境变量示例文件
COPY .env.example /app/.env

ENV PORT=5000
ENV HOST=0.0.0.0
ENV DB_PATH=school.db
ENV FLASK_DEBUG=false
ENV GIN_MODE=release

EXPOSE 5000

# 创建启动脚本
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'echo "=== 启动 school-app ==="' >> /app/start.sh && \
    echo 'echo "数据库路径: $DB_PATH"' >> /app/start.sh && \
    echo './init-db' >> /app/start.sh && \
    echo './school-app' >> /app/start.sh && \
    chmod +x /app/start.sh

CMD ["/app/start.sh"]