package grpc

import "github.com/nikitagrgv/movies/internal/delivery/grpc/pb"

type LogEntry struct {
	ClientIP string
	LastURL  string
	Usage    int64
}

type Logger interface {
	GetLogs() []LogEntry
}

type LogHandler struct {
	pb.UnimplementedLogServiceServer
	logger Logger
}

func NewLogHandler(logger Logger) *LogHandler {
	return &LogHandler{logger: logger}
}
