package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/hpcsc/background"
	"github.com/hpcsc/background/job"
)

func newHttpServer(logger *slog.Logger) *httpServer {
	return &httpServer{
		server: &http.Server{
			Addr: ":8080",
			Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				writer.Write([]byte("hello world"))
			}),
		},
		logger: logger,
	}
}

var _ job.WorkWithCleanUp = (*httpServer)(nil)

type httpServer struct {
	server *http.Server
	logger *slog.Logger
}

func (s *httpServer) CleanUp() {
	if err := s.server.Shutdown(context.Background()); err != nil {
		s.logger.Error(fmt.Sprintf("failed to stop http server: %v", err))
	}

	s.logger.Info("http server stopped")
}

func (s *httpServer) CleanUpTimeOut() time.Duration {
	return 5 * time.Second
}

func (s *httpServer) Run(logger *slog.Logger) {
	logger.Info("starting http server")

	if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Error(fmt.Sprintf("failed to start http server: %v", err))
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	runner := background.NewRunner(context.Background(), logger)

	runner.Run(job.WithRecurringFunc("job-1", 3*time.Second, func(l *slog.Logger) {
		l.Info("processing")
	}))

	runner.Run(job.WithBlockingWork("http-server", newHttpServer(logger)))

	runner.Wait()

	logger.Info("exit")
}
