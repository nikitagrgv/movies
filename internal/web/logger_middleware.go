package web

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"github.com/nikitagrgv/movies/internal/httpsrv"
	"github.com/nikitagrgv/movies/internal/logger"
)

func LoggerMiddleware(loggerService *logger.Service) httpsrv.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)

			ip, err := getUserIP(r)
			if err != nil {
				// Record at least something
				ip = netip.Addr{}
			}

			req := logger.CreateVisitRequest{
				IP:          ip,
				Path:        r.URL.Path,
				Duration:    duration,
				AttemptedAt: start,
			}
			loggerService.PushVisit(req)
		})
	}
}

func getUserIP(r *http.Request) (netip.Addr, error) {
	host := getRawHost(r)
	host = strings.TrimSpace(host)
	if ip, _, err := net.SplitHostPort(host); err == nil {
		host = ip
	}

	return netip.ParseAddr(host)
}

func getRawHost(r *http.Request) string {
	xRealIp := r.Header.Get("X-Real-IP")
	if xRealIp != "" {
		return xRealIp
	}

	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		host := xForwardedFor
		if commaIdx := strings.Index(xForwardedFor, ","); commaIdx != -1 {
			host = xForwardedFor[:commaIdx]
		}
		return host
	}

	return r.RemoteAddr
}
