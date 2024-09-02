package handlers

import (
	"context"
	"math"
	"team/transmitter/internal/domain"
	"team/transmitter/pkg/transmitter"

	"go.uber.org/zap"
)

type Messages interface {
	GetMessages(domain.AnomalyMessage)
}

type TransmitterClient struct {
	transmitterClient transmitter.TransmittersClient
	logger            *zap.Logger
	messagesServ      Messages
}

type ClientData struct {
	count    int
	sum      float64
	sum2     float64
	FillFlag bool
}

const sizeDistribution int = 100

func newClientData() *ClientData {
	return &ClientData{}
}

func NewTransmitterClient(client transmitter.TransmittersClient, logger *zap.Logger, messagesServ Messages) *TransmitterClient {
	return &TransmitterClient{
		transmitterClient: client,
		logger:            logger,
		messagesServ:      messagesServ,
	}
}

func (data *ClientData) fillData(value float64, size int) {
	data.count++
	data.sum += value
	data.sum2 += value * value

	if data.count == size {
		data.FillFlag = true
	}
}

func (data *ClientData) getMean() float64 {
	if data.count < 1 {
		return 0
	}

	return data.sum / float64(data.count)
}

func (data *ClientData) getSTD() float64 {
	if data.count < 1 {
		return 0
	}
	mean := data.getMean()

	return math.Sqrt(data.sum2/float64(data.count) - mean*mean)
}

func checkAnomaly(frequency, mean, std, kAnomaly float64) bool {
	deviation := math.Abs(frequency - mean)

	return deviation > kAnomaly*std
}

func (c *TransmitterClient) GetMessage(kAnomaly float64) {
	stream, err := c.transmitterClient.Transmit(context.Background(), &transmitter.TransmitterRequest{})
	if err != nil {
		c.logger.Fatal("Transmit failed", zap.Error(err))
	}

	data := newClientData()

	for {
		res, err := stream.Recv()
		if err != nil {
			c.logger.Fatal("Error receiving message", zap.Error(err))
		}

		if !data.FillFlag {
			data.fillData(res.Frequency, sizeDistribution)
		}

		if data.FillFlag {
			if checkAnomaly(res.Frequency, data.getMean(), data.getSTD(), kAnomaly) {
				c.logger.Info("Anomaly detected",
					zap.String("SessionID", res.SessionId),
					zap.Float64("frequency", res.Frequency),
					zap.Time("timestamp", res.Time.AsTime()))
				message := domain.AnomalyMessage{}
				message.SessionId = res.SessionId
				message.Frequency = res.Frequency
				message.Timestamp = res.Time.AsTime()
				c.messagesServ.GetMessages(message)
			}
		}
	}
}
