package main

import (
	"log"
	"net"
	"team/transmitter/internal/handlers"
	"team/transmitter/pkg/logger"
	"team/transmitter/pkg/transmitter"

	"google.golang.org/grpc"
)

func main() {
	// init logger
	logger := logger.InitLog("app.log")

	//init server
	l, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatalf("cannot create listner: %s\n", err)
	}

	serverTransmit := grpc.NewServer()
	server := handlers.NewTransmitterServer(logger)

	transmitter.RegisterTransmittersServer(serverTransmit, server)
	err = serverTransmit.Serve(l)
	if err != nil {
		log.Fatalf("impossible to serve: %s", err)
	}
}
