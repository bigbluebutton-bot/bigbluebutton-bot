module github.com/ITLab-CC/bigbluebutton-bot

go 1.19

require (
	api v1.0.0
	ddp v1.0.0
)

require (
	github.com/benpate/derp v0.22.2 // indirect
	github.com/benpate/null v0.6.4 // indirect
)

require (
	github.com/apex/log v1.9.0 // indirect
	github.com/benpate/convert v0.13.5
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pion/logging v0.2.2 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/stun v0.3.5 // indirect
	github.com/pion/transport v0.13.1 // indirect
	github.com/pion/turn/v2 v2.0.9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	golang.org/x/net v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace api v1.0.0 => ./api

replace ddp v1.0.0 => ./ddp-go
