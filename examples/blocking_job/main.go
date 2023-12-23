package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/hpcsc/background"
	"github.com/hpcsc/background/job"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	runner := background.NewRunner(context.Background(), logger)

	runner.Run(job.NewRecurring("job-1", 3*time.Second, func(_ context.Context, l *slog.Logger) error {
		l.Info("processing")
		return nil
	}))

	server := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("hello world"))
		}),
	}

	run := func(_ context.Context, l *slog.Logger) error {
		l.Info("starting http server")

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			l.Error(fmt.Sprintf("failed to start http server: %v", err))
			return err
		}

		return nil
	}

	cleanUp := func(_ context.Context, l *slog.Logger) error {
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error(fmt.Sprintf("failed to stop http server: %v", err))
			return err
		}

		logger.Info("http server stopped")
		return nil
	}
	runner.Run(job.NewBlocking("http-server", run, cleanUp))

	runner.Wait()

	logger.Info("exit")
}
