# Intergalactic Chat System

This the implementation for a Distributed System class project. The goal is to
implement a basic distributed chat system to explore recent technologies in the
field.

The whole project aims to use containerization with _Docker_, orchestration with
_Kubernetes_ and service mesh with _Istio_.

## ics-message-server

This repository refers to one of the modules of the chat system: the message
server. This service will provide an HTTP and WebSockets API to be used by a
browser client and a gRPC for andministration. It also communicates with another
service, the _name server_, that will hold information about the
_message server_ nodes state.

### API

- HTTP

|Routes               |Methods|Response           |Description                                    |
|---------------------|:-----:|:-----------------:|-----------------------------------------------|
|/auth?nickname=string|POST   |{"token": "string"}|Authenticates a new user in the message server.|

- gRPC
