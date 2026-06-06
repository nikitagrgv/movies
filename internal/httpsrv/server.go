package httpsrv

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port int, handler http.Handler) *Server {
	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return &Server{
		httpServer: srv,
	}
}

func (s *Server) Run(ctx context.Context) error {
	errChan := make(chan error, 1)
	go func() {
		fmt.Printf("HTTP Listening on address %s\n", s.httpServer.Addr)
		var ret error
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			ret = err
		}
		errChan <- ret
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		fmt.Println("\nShutting down server gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}

		return <-errChan
	}
}
