// Package postgres реализует взаимодействие с базой данных PostgreSQL.
package postgres

//nolint:depguard
import (
	"context"
	"fmt"
	"time"

	"github.com/ArtemQa8/currency-converter/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// Repository описывает репозиторий для работы с историей курсов в Postgres.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository конструктор - создаёт новый репозиторий Postgres.
func NewRepository(ctx context.Context, connString string) (*Repository, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Repository{pool: pool}, nil
}

// SaveRate сохраняет курс в базу.
func (r *Repository) SaveRate(ctx context.Context, rate model.Rate) error {
	sqlQuery := `
		INSERT INTO rate_history (base_currency, quote_currency, value, date, source)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(ctx, sqlQuery,
		rate.BaseCurrency,
		rate.QuoteCurrency,
		rate.Value.String(),
		rate.Date,
		rate.Source,
	)

	return err
}

// GetRate получает последний курс из базы.
func (r *Repository) GetRate(ctx context.Context, baseCurrency, quoteCurrency string) (model.Rate, error) {
	sqlQuery := `
		SELECT base_currency, quote_currency, value, date, source
		FROM rate_history
		WHERE base_currency = $1 AND quote_currency = $2
		ORDER BY date DESC
		LIMIT 1
	`

	var base, quote, source, valStr string
	var date time.Time

	err := r.pool.QueryRow(ctx, sqlQuery, baseCurrency, quoteCurrency).
		Scan(&base, &quote, &valStr, &date, &source)
	if err != nil {
		return model.Rate{}, err
	}

	valDec, err := decimal.NewFromString(valStr)
	if err != nil {
		return model.Rate{}, err
	}

	return model.Rate{
		BaseCurrency:  base,
		QuoteCurrency: quote,
		Value:         valDec,
		Date:          date,
		Source:        "postgres",
	}, nil
}

// GetHistory возвращает историю курсов за указанное количество дней.
func (r *Repository) GetHistory(ctx context.Context,
	baseCurrency, quoteCurrency string, days int,
) ([]model.Rate, error) {
	query := `
		SELECT base_currency, quote_currency, value, date, source
		FROM rate_history
		WHERE base_currency = $1 AND quote_currency = $2
			AND date >= NOW() - ($3 * INTERVAL '1 day')
		ORDER BY date ASC
	`
	rows, err := r.pool.Query(ctx, query, baseCurrency, quoteCurrency, days)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса истории: %w", err)
	}
	defer rows.Close()

	var rates []model.Rate
	for rows.Next() {
		var base, quote, source, valStr string
		var date time.Time

		if err := rows.Scan(&base, &quote, &valStr, &date, &source); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}

		valDec, err := decimal.NewFromString(valStr)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга decimal: %w", err)
		}

		rates = append(rates, model.Rate{
			BaseCurrency:  base,
			QuoteCurrency: quote,
			Value:         valDec,
			Date:          date,
			Source:        source,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам: %w", err)
	}

	return rates, nil
}
