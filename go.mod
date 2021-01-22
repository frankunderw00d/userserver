module userserver

go 1.14

require (
	github.com/gomodule/redigo v1.8.3
	jarvis v1.0.1
	baseservice v1.0.1
)

replace (
    jarvis v1.0.1 => /home/frank/Documents/project/jarvis
    baseservice v1.0.1 => /home/frank/Documents/project/baseservice
)
