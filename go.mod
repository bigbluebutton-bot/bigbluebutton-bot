module github.com/ITLab-CC/bigbluebutton-bot

go 1.19

require api v1.0.0

require (
	github.com/apex/log v1.9.0 // indirect
	github.com/gopackage/ddp v0.0.4 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
)

replace api v1.0.0 => ./api
