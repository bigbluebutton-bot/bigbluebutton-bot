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
	github.com/go-rel/changeset v1.2.0
	github.com/gorilla/websocket v1.5.0
	github.com/graarh/golang-socketio v1.0.0
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d
)

require (
	github.com/apex/log v1.9.0 // indirect
	github.com/azer/snakecase v0.0.0-20161028114325-c818dddafb5c // indirect
	github.com/benpate/derp v0.22.2 // indirect
	github.com/benpate/null v0.6.4 // indirect
	github.com/go-rel/rel v0.30.0 // indirect
	github.com/gopackage/ddp v0.0.5 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/robertkrimen/otto v0.2.1 // indirect
	github.com/serenize/snaker v0.0.0-20201027110005-a7ad2135616e // indirect
	github.com/tidwall/gjson v1.11.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
)
