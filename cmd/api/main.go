package main

import "github.com/abyan-dev/auth/pkg/server"

func main() {
	app := server.New()
	server.Run(app)
}
