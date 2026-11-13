package postgres

import (
	"context"
	"time"

	"github.com/ArtemQa8/currency-converter/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Repository struct {
	pool *pgxpool.Pool
}

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
		Source:        source,
	}, nil
}
