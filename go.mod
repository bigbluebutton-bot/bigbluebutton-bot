module github.com/ITLab-CC/bigbluebutton-bot

go 1.19

replace api v1.0.0 => ./api
replace ddp v1.0.0 => ./ddp-go
replace go-socket.io-client v1.0.0 => ./go-socket.io-client //old
replace github.com/graarh/golang-socketio v1.0.0 => ./golang-socketio

require (
	api v1.0.0
	ddp v1.0.0
	github.com/benpate/convert v0.13.5
	github.com/gorilla/websocket v1.5.0
	github.com/graarh/golang-socketio v1.0.0
)

require (
	github.com/apex/log v1.9.0 // indirect
	github.com/benpate/derp v0.22.2 // indirect
	github.com/benpate/null v0.6.4 // indirect
	github.com/gopackage/ddp v0.0.5 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
)
