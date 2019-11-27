FROM golang:alpine AS build-env

RUN apk --no-cache add build-base git gcc
COPY go.mod main.go server.go /src/
COPY api /src/api
COPY server /src/server

RUN cd /src && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ics-message-server

FROM scratch

COPY --from=build-env /src/main /

ENV ICS_JWT_TOKEN="distributed_systems_rules"

EXPOSE 7000
