# ics-message-server

This repository refers to one of the modules of a chat system: the message
server. This service will provide an HTTP and WebSockets API to be used by a
browser client and a gRPC for andministration. It also communicates with another
service, the _name server_, that will hold information about the
_message server_ nodes state.

## Usage

### Build

```sh
git clone https://github.com/renantarouco/ics-message-server.git
cd ics-message-server
go test -v ./...
go build
```

### Run

After building you can run the executable.

```sh
./ics-message-server
```

If you chose to run `go install` command than you can head directly to your
_GOPATH_ directory and run the executable, or if you have it in your system's
_PATH_ then you just have to run `ics-message-server`.

## API

- HTTP

|Routes               |Methods|Response           |Description                                    |
|---------------------|:-----:|:-----------------:|-----------------------------------------------|
|/auth?nickname=string|POST   |{"token": "string"}|Authenticates a new user in the message server.|

- gRPC

  TODO

## Intergalactic Chat System

This the implementation for a Distributed System class project. The goal is to
implement a basic distributed chat system to explore recent technologies in the
field.

The whole project aims to use containerization with _Docker_, orchestration with
_Kubernetes_ and service mesh with _Istio_.
