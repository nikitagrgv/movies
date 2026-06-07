package logger

import (
	"net/netip"
	"time"
)

type CreateVisitRequest struct {
	IP          netip.Addr
	Path        string
	Duration    time.Duration
	AttemptedAt time.Time
}

type Visit struct {
	ID          int
	IP          netip.Addr
	Path        string
	Duration    time.Duration
	AttemptedAt time.Time
}
