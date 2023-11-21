package main

import (
	"context"
	"fmt"
	"github.com/hpcsc/background"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var _ background.Job = (*httpServerJob)(nil)

type httpServerJob struct{}

func (j *httpServerJob) Name() string {
	return "http-server"
}

func (j *httpServerJob) Run(ctx context.Context, logger *slog.Logger) {
	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("hello world"))
		}),
	}

	go func() {
		logger.Info("starting http server")

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("failed to start http server: %v", err))
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Error(fmt.Sprintf("failed to stop http server: %v", err))
	} else {
		logger.Info("http server stopped")
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	runner := background.NewRunner(context.Background(), logger)

	runner.Run(&httpServerJob{})
	runner.Run(background.NewTickerJob("ticker-job", 5*time.Second, func(_ context.Context, l *slog.Logger) error {
		l.Info("processing")
		return nil
	}))

	runner.Wait()

	logger.Info("exit")
}
