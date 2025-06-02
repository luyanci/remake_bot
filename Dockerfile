FROM golang:latest AS build
 
WORKDIR /src
 
COPY . .
 
RUN go mod tidy -v;CGO_ENABLED=0 go build -o /bot

FROM alpine:latest
WORKDIR /src

COPY --from=build /bot bot
COPY --from=build /src/countries.json countries.json
COPY --from=build /src/user_list.json user_list.json

CMD [ "/src/bot" ]
