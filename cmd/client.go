package cmd

import (
	"context"
	"fmt"
	"github.com/bcessa/sample-grpc-project/proto"
	"github.com/bcessa/sample-grpc-project/rpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"path"
	"time"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start sample client",
	RunE:  runClient,
}

func init() {
	var (
		unix    bool
		port    int
		timeout int
	)
	clientCmd.Flags().BoolVar(&unix, "unix", false, "use unix socket")
	clientCmd.Flags().IntVar(&port, "port", 9000, "tcp port to use")
	clientCmd.Flags().IntVar(&timeout, "timeout", 5, "connection timeout in seconds")
	viper.BindPFlag("client.unix", clientCmd.Flags().Lookup("unix"))
	viper.BindPFlag("client.port", clientCmd.Flags().Lookup("port"))
	viper.BindPFlag("client.timeout", clientCmd.Flags().Lookup("timeout"))
	RootCmd.AddCommand(clientCmd)
}

func runClient(_ *cobra.Command, _ []string) error {
	// Use context to set a timeout when dialing the gRPC connection;
	// this is required since 'WithTimeout()' has been deprecated
	timeout := time.Duration(viper.GetInt("client.timeout")) * time.Second
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	// Basic connection options
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}

	var location string
	if viper.GetBool("client.unix") {
		// Custom dialer is required to use a unix socket
		location = path.Join(".", "sample-server.sock")
		opts = append(opts, grpc.WithDialer(func(address string, _ time.Duration) (net.Conn, error) {
			return net.Dial("unix", address)
		}))
	} else {
		location = fmt.Sprintf(":%d", viper.GetInt("client.port"))
	}

	// Get connection
	fmt.Printf("requesting network connection at: %s\n", location)
	conn, err := grpc.DialContext(ctx, location, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()
	fmt.Printf("connection state: %s\n", conn.GetState())

	// Get client
	client := proto.NewSampleServiceClient(conn)

	// Start interactive console for commands
	fmt.Println("starting interactive commands console")
	cli := rpc.NewConsole(client, "\033[33m»\033[0m ")
	defer cli.Close()
	return cli.Start()
}
