package main

//nolint:depguard
import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArtemQa8/currency-converter/internal/api"
	"github.com/ArtemQa8/currency-converter/internal/build"
	"github.com/ArtemQa8/currency-converter/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	ctx := context.Background()
	svc, err := build.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Ошибка сборки приложения: %v", err)
	}

	handler := api.NewHandler(svc)
	router := api.NewRouter(handler)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("Сервер запущен на порту :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	log.Println("Получен сигнал остановки, завершаем работу...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при graceful shutdown: %v", err)
	}

	log.Println("Сервер остановлен.")
}
