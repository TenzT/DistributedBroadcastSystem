package main

import (
	"BroadcastService/server"
)

func main() {
	svc := server.New()
	svc.Run()
}
