# rubbit

Rubbit is a sample project i did to get more familiar with RabbitMQ and its plugins(delayed message).

# Instalation
To install the Rubbit package use the command below:

go get github.com/bamdadam/rubbit

# Instructions

## Run

* run services and server: `make run-server`
* run services: `make up`
* run clients: `go run main.go client {topic names}`

## Example

### Posting Event
*request
```
curl -X POST http://127.0.0.1:8080/publish -H 'Content-Type: application/json' -d '{"topic":"test-topic", "message": "first-message", "publish_delay":"5000ms", "delayed":true}'
```