FROM golang:latest AS build
 
WORKDIR /src
 
COPY . .
 
RUN go mod tidy -v;CGO_ENABLED=0 go build -o /bot

FROM alpine:latest
WORKDIR /src
RUN apk update && \
    apk add --no-cache postgresql postgresql-client && \
    apk add --no-cache build-base && \
    rm -rf /var/cache/apk/*

RUN mkdir -p /var/lib/postgresql/data
RUN chown -R postgres:postgres /var/lib/postgresql/data

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=114514
ENV POSTGRES_DB=postgres

RUN su - postgres -c "initdb -D /var/lib/postgresql/data"

COPY start_postgres.sh /start_postgres.sh
RUN chmod +x /start_postgres.sh

COPY --from=build /bot bot
COPY --from=build /src/countries.json countries.json
COPY --from=build /src/user_list.json user_list.json

EXPOSE 5432
CMD ["/start_postgres.sh", "/src/bot" ]
