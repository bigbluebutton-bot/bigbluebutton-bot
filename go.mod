module github.com/ITLab-CC/bigbluebutton-bot

go 1.19

require (
	api v1.0.0
)

replace (
	api v1.0.0 => ./api
)