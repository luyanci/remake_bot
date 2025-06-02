#!/bin/sh
set -e

# 启动 PostgreSQL 服务
su - postgres -c "pg_ctl -D /var/lib/postgresql/data -l logfile start"

# 等待数据库启动（可以根据实际情况调整等待时间或检查方式）
sleep 5

# 运行传入的命令（这里是 /src/bot）
exec "$@"
