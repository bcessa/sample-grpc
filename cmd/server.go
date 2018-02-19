package cmd

import (
	"context"
	"fmt"
	"github.com/bcessa/sample-grpc/proto"
	"github.com/bcessa/sample-grpc/rpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start sample server",
	RunE:  runServer,
}

func init() {
	var (
		unix    bool
		port    int
		useREST bool
	)
	serverCmd.Flags().BoolVar(&unix, "unix", false, "use unix socket")
	serverCmd.Flags().BoolVar(&useREST, "rest", true, "enable REST interface")
	serverCmd.Flags().IntVar(&port, "port", 9000, "tcp port to use")
	viper.BindPFlag("server.unix", serverCmd.Flags().Lookup("unix"))
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
	viper.BindPFlag("server.rest", serverCmd.Flags().Lookup("rest"))
	RootCmd.AddCommand(serverCmd)
}

func runServer(_ *cobra.Command, _ []string) error {
	var l net.Listener
	var location string
	var err error

	if viper.GetBool("server.unix") {
		// Get socket location with proper cleanup
		location = path.Join(".", "sample-server.sock")
		defer os.Remove(location)

		// Get network listener for socket file
		l, err = net.Listen("unix", path.Join(".", "sample-server.sock"))
		if err != nil {
			return err
		}
	} else {
		location = fmt.Sprintf(":%d", viper.GetInt("server.port"))
		l, err = net.Listen("tcp", location)
		if err != nil {
			return err
		}
	}

	// Register sample service handler
	server := grpc.NewServer()
	proto.RegisterSampleServiceServer(server, &rpc.SampleService{})
	reflection.Register(server)

	// Start service processing on a different thread (go routine)
	log.Printf("starting server on: %s", location)
	go server.Serve(l)

	// HTTP gateway
	if viper.GetBool("server.rest") && !viper.GetBool("server.unix") {
		restLoc := fmt.Sprintf(":%d", viper.GetInt("server.port") + 1)
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithInsecure()}
		err := proto.RegisterSampleServiceHandlerFromEndpoint(context.TODO(), mux, location, opts)
		if err != nil {
			return err
		}

		log.Printf("starting REST interface %s", restLoc)
		go http.ListenAndServe(restLoc, mux)
	}

	// Wait for system signals
	signalsCh := make(chan os.Signal)
	signal.Notify(signalsCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		select {
		case <-signalsCh:
			log.Println("gracefully stopping server")
			server.GracefulStop()
			return nil
		}
	}
}
