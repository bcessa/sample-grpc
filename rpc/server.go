package rpc

import (
	"context"
	"github.com/bcessa/sample-grpc/proto"
	"github.com/gogo/protobuf/types"
	"time"
)

type SampleService struct{}

func (s *SampleService) Ping(_ context.Context, _ *types.Empty) (*proto.Pong, error) {
	return &proto.Pong{Ok: true}, nil
}

func (s *SampleService) Items(_ *types.Empty, stream proto.SampleService_ItemsServer) error {
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
