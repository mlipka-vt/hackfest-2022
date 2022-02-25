package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"kafka-go-client/trade"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

var trades sync.Map

type Service struct {
	trade.UnimplementedTradeServiceServer
}

func (s Service) GetTradeStatus(ctx context.Context, request *trade.GetTradeStatusRequest) (*trade.GetTradeStatusResponse, error) {
	fmt.Println(fmt.Sprintf("INFO received trade status request for %s", request.Id))
	t, ok := trades.Load(request.Id)
	if !ok {
		fmt.Println("ERROR invalid uuid provided for trade status")
		return nil, fmt.Errorf("invalid uuid")
	}

	return &trade.GetTradeStatusResponse{Trade: t.(*trade.Trade)}, nil
}

func main() {
	trades = sync.Map{}

	go tradeListener()
	go ticker()

	lis, err := net.Listen("tcp", fmt.Sprintf(":50052"))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	trade.RegisterTradeServiceServer(s, &Service{})

	fmt.Println("INFO trade-service listening")
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func ticker() {
	rand.Seed(time.Now().Unix())
	for {

		value := rand.Float64() * 100
		trades.Range(func(key, val interface{}) bool {
			t := val.(*trade.Trade)
			if t.Amount >= value && t.Status != "filled" {
				fmt.Println(fmt.Sprintf("INFO filling trade %s", t.Id))
				t.Status = "filled"
				return false
			}
			return true
		})
	}
}

func tradeListener() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "trades",
		GroupID:   "trade-receiver",
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			panic(err)
		}

		var message trade.Trade
		err = proto.Unmarshal(m.Value, &message)
		if err != nil {
			log.Println(fmt.Errorf("ERROR unmarshalling trade message: %w"))
			continue
		}

		fmt.Println(fmt.Sprintf("INFO received trade: %s", message.String()))
		trades.Store(message.Id, &message)
	}
}
