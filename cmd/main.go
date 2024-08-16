package main

// Логирование в проекте не делал для экономии времени, но, естественно, оно должно быть)

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	orderHttp "hotel/internal/http"
	"hotel/internal/repository"
	"hotel/internal/service"
)

func main() {
	orderRepo := repository.NewOrderRepository()
	availabilityRepo := repository.NewAvailabilityRepository()
	srv := service.New(orderRepo, availabilityRepo)
	h := orderHttp.New(srv)

	mux := http.NewServeMux()
	mux.HandleFunc("/orders", h.CreateOrder)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		<-stop

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			LogErrorf("Server forced to shutdown: %v", err)
		}

		close(done)
	}()

	go func() {
		LogInfo("Server listening on localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			LogErrorf("Server failed: %v", err)
		}
	}()

	<-done
	wg.Wait()
	LogInfo("Server stopped gracefully")
}

var logger = log.Default()

func LogErrorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Error]: %s\n", msg)
}

func LogInfo(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	logger.Printf("[Info]: %s\n", msg)
}
