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
	port   int
	server *grpc.Server
}

func NewServer(port int) *Server {
	return &Server{
		port:   port,
		server: grpc.NewServer(),
	}
}

func (s *Server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(s.port))
	if err != nil {
		return err
	}
	fmt.Printf("gRPC Listening on address %s\n", lis.Addr().String())

	defer s.server.GracefulStop()

	errChan := make(chan error, 1)
	go func() {
		var ret error
		if err := s.server.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
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
