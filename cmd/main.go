package main

import (
	"BroadcastService/server"
	"flag"
	"log"
	"strconv"
)

func main() {

	httpPort := flag.String("httpport", "", "http port for clients.")

	flag.Parse()

	if *httpPort == "" {
		log.Println("empty httpport.")
		return
	}
	port, err := strconv.Atoi(*httpPort)
	if err != nil {
		log.Println(err.Error())
		return
	}
	serverConfig := &server.ServerConfig{
		HttpPort: int64(port),
	}
	svc := server.New(serverConfig)
	svc.Run()
}
