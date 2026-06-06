package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type Server struct {
	port int
}

func NewServer(port int) *Server {
	return &Server{port: port}
}

func (s *Server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(s.port))
	if err != nil {
		return err
	}
	fmt.Printf("gRPC Listening on address %s\n", lis.Addr().String())

	server := grpc.NewServer()
	defer server.GracefulStop()

	errChan := make(chan error, 1)
	go func() {
		var ret error
		if err := server.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			ret = fmt.Errorf("gRPC server error: %w", err)
		}
		errChan <- ret
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		fmt.Println("\nShutting down gRPC server gracefully...")
		return ctx.Err()
	}
}
