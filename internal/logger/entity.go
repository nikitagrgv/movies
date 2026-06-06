package logger

import (
	"net/netip"
	"time"
)

type CreateVisitRequest struct {
	IP          netip.Addr
	AttemptedAt time.Time
	Path        string
}

type Visit struct {
	ID          int
	IP          netip.Addr
	AttemptedAt time.Time
	Path        string
}
