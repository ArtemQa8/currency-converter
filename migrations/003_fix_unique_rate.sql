-- +goose Up
DROP INDEX IF EXISTS idx_rate_history_unique;
CREATE UNIQUE INDEX idx_rate_history_unique ON rate_history (base_currency, quote_currency, date, value);

-- +goose Down
DROP INDEX idx_rate_history_unique;
CREATE UNIQUE INDEX idx_rate_history_unique ON rate_history (base_currency, quote_currency, date);