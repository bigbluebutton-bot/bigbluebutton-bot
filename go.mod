module github.com/ITLab-CC/bigbluebutton-bot

go 1.19

require (
	api v1.0.0
	ddp v1.0.0
)

require (
	github.com/benpate/derp v0.22.2 // indirect
	github.com/benpate/null v0.6.4 // indirect
	github.com/gopackage/ddp v0.0.5 // indirect
)

require (
	github.com/apex/log v1.9.0 // indirect
	github.com/benpate/convert v0.13.5
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	golang.org/x/net v0.7.0 // indirect
)

replace api v1.0.0 => ./api

replace ddp v1.0.0 => ./ddp-go
