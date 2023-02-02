module github.com/ITLab-CC/bigbluebutton-bot

go 1.19

require (
	api v1.0.0
	ddp v1.0.0
	github.com/zhouhui8915/engine.io-go v0.0.0-20150910083302-02ea08f0971f
	go-socket.io-client v1.0.0
)

require (
	github.com/benpate/derp v0.22.2 // indirect
	github.com/benpate/null v0.6.4 // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/gomodule/redigo v1.8.4 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
)

require (
	github.com/apex/log v1.9.0 // indirect
	github.com/benpate/convert v0.13.5
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/googollee/go-socket.io v1.6.2
	github.com/pion/logging v0.2.2 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/stun v0.3.5 // indirect
	github.com/pion/transport v0.13.1 // indirect
	github.com/pion/turn/v2 v2.0.9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	github.com/zhouhui8915/go-socket.io-client v0.0.0-20200925034401-83ee73793ba4
	golang.org/x/net v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace api v1.0.0 => ./api

replace ddp v1.0.0 => ./ddp-go

replace go-socket.io-client v1.0.0 => ./go-socket.io-client
