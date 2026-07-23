module github.com/RAN-GAN/rendis/client/cli

go 1.26.4

require (
	github.com/RAN-GAN/rendis/client/golang v0.0.0
	github.com/chzyer/readline v1.5.1
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/sys v0.0.0-20220310020820-b874c991c1a5 // indirect
)

replace github.com/RAN-GAN/rendis/client/golang => ../golang
