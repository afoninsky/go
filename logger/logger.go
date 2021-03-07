package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

var nilTraceID trace.TraceID

// Logger ...
type Logger struct {
	logrus.Entry
}

// New creates new logger instance
func New() *Logger {
	var log = logrus.NewEntry(logrus.New())

	log.Logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	return &Logger{*log}
}

// WithContext creates new logger with http request context
func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
	logger := l.Dup()

	// add trace id from the context if it exists
	span := trace.SpanFromContext(ctx)
	traceID := span.SpanContext().TraceID
	if traceID != nilTraceID {
		logger = l.WithFields(logrus.Fields{
			"trace": span.SpanContext().TraceID,
		})
	}

	return logger
}

// Middleware returns http logging middleware
func (l *Logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(w, r)

		l.WithContext(r.Context()).
			WithField("method", r.Method).
			WithField("status", wrapped.status).
			// WithField("path", r.URL.EscapedPath()).
			WithField("duration", time.Since(start)).
			Info(r.URL)
	})
}
