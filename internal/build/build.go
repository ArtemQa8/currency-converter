package build

//nolint:depguard
import (
	"context"
	"fmt"

	"github.com/ArtemQa8/currency-converter/internal/config"
	"github.com/ArtemQa8/currency-converter/internal/repository/cbr"
	"github.com/ArtemQa8/currency-converter/internal/repository/file"
	"github.com/ArtemQa8/currency-converter/internal/repository/postgres"
	"github.com/ArtemQa8/currency-converter/internal/repository/redis"
	"github.com/ArtemQa8/currency-converter/internal/service"
)

// New собирает приложение: подключает базы, создает клиенты и инициализирует сервис
func New(ctx context.Context, cfg *config.Config) (*service.ConverterService, error) {
	pgRepo, err := postgres.NewRepository(ctx, cfg.PostgresURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к Postgres: %w", err)
	}

	redisCache := redis.NewCache(cfg.RedisURL, cfg.RedisTTL)
	cbrClient := cbr.NewClient(cfg.CBRAPIURL)
	fileRepo := file.NewRepository(cfg.RatesFilePath)

	// Порядок важен (цепочка отказоустойчивости)
	svc := service.NewConverterService(
		redisCache, // Сначала кэш
		cbrClient,  // ЦБ РФ
		pgRepo,     // История из базы
		fileRepo,   // Файл
	)

	return svc, nil
}
