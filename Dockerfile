FROM golang:alpine AS build-env

RUN apk --no-cache add build-base git gcc
COPY go.mod main.go server.go /src/
COPY api /src/api
COPY server /src/server

RUN cd /src && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main

FROM scratch

COPY --from=build-env /src/main /

EXPOSE 7000

ENTRYPOINT [ "/main" ]