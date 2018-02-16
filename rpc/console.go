package rpc

import (
	"context"
	"fmt"
	"github.com/bcessa/sample-grpc/proto"
	"github.com/chzyer/readline"
	"io"
	"log"
	"time"
)

type ClientConsole struct {
	client proto.SampleServiceClient
	rl     *readline.Instance
}

func NewConsole(c proto.SampleServiceClient, prompt string) *ClientConsole {
	rl, _ := readline.New(prompt)
	return &ClientConsole{
		client: c,
		rl:     rl,
	}
}

func (c *ClientConsole) Start() error {
	c.usage()
	for {
		line, err := c.rl.Readline()
		if err != nil {
			return err
		}
		switch line {
		case "p":
			pong, _ := c.client.Ping(context.TODO(), &proto.Empty{})
			log.Printf("pong: %v\n", pong.Ok)
		case "s":
			ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
			stream, err := c.client.Items(ctx, &proto.Empty{})
			if err != nil {
				return err
			}
			for {
				item, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Println(err)
					break
				}
				log.Printf("item: %d\n", item.Id)
			}
			fmt.Println("finish stream processing")
			cancel()
		case "q":
			fmt.Println("closing console")
			return nil
		case "h":
			c.usage()
		default:
			fmt.Println("invalid command")
		}
	}
}

func (c *ClientConsole) Close() {
	c.rl.Close()
}

func (c *ClientConsole) usage() {
	fmt.Println("p = ping")
	fmt.Println("s = stream")
	fmt.Println("h = help")
	fmt.Println("q = quit")
}
