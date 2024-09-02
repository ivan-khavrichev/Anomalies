package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"team/transmitter/internal/handlers"
	"team/transmitter/internal/repository"
	service "team/transmitter/internal/services"
	psql "team/transmitter/pkg/database"
	"team/transmitter/pkg/logger"
	"team/transmitter/pkg/transmitter"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// flag deviation
	var kFlag = flag.Float64("k", 3, "deviation coefficient")
	flag.Parse()

	// init loging
	logger := logger.InitLog("app.log")

	// init dotenv
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("cant load .env", err)
	}

	// read config to db
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", os.Getenv("host"), os.Getenv("username"), os.Getenv("password"), os.Getenv("dbname"), os.Getenv("port"))

	// init connetcion to database
	db, err := psql.ConnectDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	// connect to server
	conn, err := grpc.NewClient("localhost:8089", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// inital client and connetcion
	client := transmitter.NewTransmittersClient(conn)
	rep := repository.NewRepository(db)
	serv := service.NewMessages(rep)
	trClient := handlers.NewTransmitterClient(client, logger, serv)
	trClient.GetMessage(*kFlag)
}
