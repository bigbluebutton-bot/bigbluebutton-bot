module github.com/ITLab-CC/bigbluebutton-bot

go 1.19

require (
	api v1.0.0
	ddp v1.0.0
)

require (
	github.com/apex/log v1.9.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	golang.org/x/net v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace api v1.0.0 => ./api

replace ddp v1.0.0 => ./ddp-go
