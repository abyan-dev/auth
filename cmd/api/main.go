package main

import "github.com/abyan-dev/auth/pkg/server"

func main() {
	srv := server.Server{}
	app := srv.New()
	srv.Run(app)
}
