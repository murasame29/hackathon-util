package middleware

import (
	"net/http"
	"time"

	"github.com/murasame29/hackathon-util/pkg/logger"
)

func BuildChain(h http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
	if len(m) == 0 {
		return h
	}
	return m[0](BuildChain(h, m[1:cap(m)]...))
}

func LoggerInContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.NewLoggerWithContext(r.Context())

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

type accessLog struct {
	Agent         string        `json:"user_agent"`
	Referer       string        `json:"referer"`
	RemoteIP      string        `json:"remote_ip"`
	RequestID     string        `json:"request_id"`
	RequestMethod string        `hson:"request_method"`
	RequestURI    string        `json:"request_uri"`
	TimeStamp     time.Time     `json:"time_stamp"`
	Latency       time.Duration `json:"latency(ms)"`
}

func AccessLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var accessLog accessLog
		accessLog.Agent = r.UserAgent()
		accessLog.Referer = r.Referer()
		accessLog.RemoteIP = r.RemoteAddr
		accessLog.RequestID = r.Header.Get("X-Request-ID")
		accessLog.RequestMethod = r.Method
		accessLog.RequestURI = r.RequestURI
		accessLog.TimeStamp = time.Now()

		h.ServeHTTP(w, r)
		accessLog.Latency = time.Since(accessLog.TimeStamp) / time.Millisecond
		logger.Info(r.Context(), "access log", logger.Field("accessLog", accessLog))
	})
}
