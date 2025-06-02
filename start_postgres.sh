#!/bin/bash
set -e

# 等待 PostgreSQL 服务启动并就绪
echo "Waiting for PostgreSQL to start..."
while! pg_isready -U $POSTGRES_USER -d $POSTGRES_DB -h localhost -p 5432; do
    sleep 1
done
echo "PostgreSQL is ready."

# 在这里添加启动你应用程序的命令，假设你的应用程序可执行文件名为 app
# 请根据实际情况替换为你的应用程序启动命令
/bot

# 如果你的应用程序是在后台运行的，可以使用以下方式启动并等待
# /app &
# wait $!
