package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/config" //nolint:depguard
)

type Server struct {
	logger Logger
	app    Application
	server *http.Server
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, id, title string) error
}

func NewServer(logger Logger, app Application, conf config.ServerConf) *Server {
	mux := http.NewServeMux()
	mux.Handle("/hello", logMiddleware(http.HandlerFunc(helloHandler), logger))

	return &Server{
		logger: logger,
		app:    app,
		server: &http.Server{
			Addr:              conf.Host + ":" + strconv.Itoa(conf.Port),
			Handler:           mux,
			ReadHeaderTimeout: 10 * time.Second, // Adjust the timeout as necessary
		},
	}
}

func logMiddleware(next http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isUnwantedRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		clientIP := r.RemoteAddr
		userAgent := r.UserAgent()
		method := r.Method
		path := r.URL.Path
		proto := r.Proto

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		latency := time.Since(start)
		logger.Info(fmt.Sprintf("%s [%s] %s %s %s %d %d \"%s\"", clientIP,
			start.Format(time.RFC1123), method, path, proto, rec.status, latency.Milliseconds(), userAgent))
	})
}

func isUnwantedRequest(r *http.Request) bool {
	unwantedPaths := []string{"/favicon.ico", "/static/"}
	for _, prefix := range unwantedPaths {
		if strings.HasPrefix(r.URL.Path, prefix) {
			return true
		}
	}

	// Skip logging for 404 pages
	return r.Method == http.MethodGet && r.URL.Path == "/non-existent"
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("hello Otus Student!"))
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Stop(ctxShutdown); err != nil {
			s.logger.Error(fmt.Sprintf("Server shutdown error: %v", err))
		}
	}()

	s.logger.Info(fmt.Sprintf("Starting server on %s", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("Server was stopped on %s", s.server.Addr))
	return nil
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}
