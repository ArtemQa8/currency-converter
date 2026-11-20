// Package main запускает HTTP-сервер.
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

	ctx, cancel := context.WithCancel(context.Background())

	svc, err := build.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Ошибка сборки приложения: %v", err)
	}

	defer cancel()

	handler := api.NewHandler(svc)
	router := api.NewRouter(handler)

	log.Println("Первоначальная загрузка курсов...")
	if err := svc.UpdateAllRates(ctx); err != nil {
		log.Printf("Ошибка первичной загрузки: %v", err)
	}

	go func() {
		ticker := time.NewTicker(3 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.Println("Фоновое обновление курсов...")
				if err := svc.UpdateAllRates(ctx); err != nil {
					log.Printf("Ошибка фонового обновления: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

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

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Ошибка при graceful shutdown: %v", err)
	}

	log.Println("Сервер остановлен.")
}
