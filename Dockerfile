FROM golang:latest AS build
 
WORKDIR /src
 
COPY . .
 
RUN go mod tidy -v;CGO_ENABLED=0 go build -o /bot

# 使用官方的 PostgreSQL 基础镜像
FROM postgres:14

RUN apk update && apk add --no-cache bash

# 设置环境变量，定义数据库用户名和密码
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=114514
ENV POSTGRES_DB=postgres

# 将初始化 SQL 脚本复制到镜像中（如果有额外的初始化需求）
# COPY init.sql /docker-entrypoint-initdb.d/
COPY --from=build /bot /bot
COPY --from=build /src/countries.json /countries.json
COPY --from=build /src/user_list.json /user_list.json


# 暴露 PostgreSQL 的默认端口
EXPOSE 5432

# 复制启动脚本到镜像中
COPY start_postgres.sh /start.sh
RUN chmod +x /start.sh

# 定义容器启动时执行的命令，运行启动脚本
CMD ["/start.sh"]