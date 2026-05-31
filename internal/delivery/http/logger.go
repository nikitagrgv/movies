package http

import (
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/nikitagrgv/movies/internal/pkg/cache"
)

type logInfo struct {
	count   int64
	lastURL atomic.Pointer[string]
}

func (log *logInfo) addUsage() {
	atomic.AddInt64(&log.count, 1)
}

func (log *logInfo) getUsage() int64 {
	return atomic.LoadInt64(&log.count)
}

func (log *logInfo) setLastURL(url string) {
	log.lastURL.Store(&url)
}

func (log *logInfo) getLastURL() string {
	p := log.lastURL.Load()
	if p == nil {
		return ""
	}
	return *p
}

type LoggerMiddleware struct {
	data *cache.LRUCache[string, *logInfo]
}

func NewLoggerMiddleware(size int) *LoggerMiddleware {
	c := cache.NewLRUCache[string, *logInfo](size)
	return &LoggerMiddleware{
		data: c,
	}
}

func (m *LoggerMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		data, _ := m.data.GetOrPut(ip, func() *logInfo {
			return &logInfo{}
		})

		data.addUsage()
		data.setLastURL(r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

func getClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff == "" {
		return r.RemoteAddr
	}

	parts := strings.Split(xff, ",")
	return strings.TrimSpace(parts[0])
}
