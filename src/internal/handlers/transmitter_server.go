package handlers

import (
	"math/rand"
	"team/transmitter/pkg/transmitter"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransmitterServer struct {
	transmitter.UnimplementedTransmittersServer
	logger *zap.Logger
}

func NewTransmitterServer(logger *zap.Logger) *TransmitterServer {
	return &TransmitterServer{
		logger: logger,
	}
}

func (serv *TransmitterServer) Transmit(
	req *transmitter.TransmitterRequest,
	stream transmitter.Transmitters_TransmitServer) error {

	sessionID := uuid.New().String()
	mean := rand.Float64()*20 - 10
	std := rand.Float64()*1.2 + 0.3

	serv.logger.Info("New Connection", zap.String("SessionID", sessionID), zap.Float64("mean", mean), zap.Float64("std", std))
	res := &transmitter.TransmitterResponse{
		SessionId: sessionID,
	}

	for {
		res.Frequency = rand.NormFloat64()*std + mean
		res.Time = timestamppb.Now()

		err := stream.Send(res)
		if err != nil {
			return err
		}
		serv.logger.Info("Connection info", zap.String("SessionID", res.SessionId), zap.Float64("frequency", res.Frequency), zap.Time("timestamp", res.Time.AsTime()))
		time.Sleep(time.Millisecond)
	}
}
