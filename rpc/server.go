package rpc

import (
	"context"
	"github.com/bcessa/sample-grpc-project/proto"
	"time"
)

type SampleService struct{}

func (s *SampleService) Ping(_ context.Context, _ *proto.Empty) (*proto.Pong, error) {
	return &proto.Pong{Ok: true}, nil
}

func (s *SampleService) Items(_ *proto.Empty, stream proto.SampleService_ItemsServer) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	counter := 0
	for {
		select {
		case <-ticker.C:
			if err := stream.Send(&proto.Item{Id: int32(counter)}); err != nil {
				return err
			}
			if counter >= 500 {
				return nil
			}
			counter++
		case <-stream.Context().Done():
			return nil
		}
	}
}
