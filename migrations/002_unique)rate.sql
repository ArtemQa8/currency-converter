-- +goose Up
DELETE FROM rate_history a USING rate_history b
WHERE a.id > b.id 
  AND a.base_currency = b.base_currency 
  AND a.quote_currency = b.quote_currency 
  AND a.date = b.date;

CREATE UNIQUE INDEX idx_rate_history_unique ON rate_history (base_currency, quote_currency, date);

-- +goose Down
DROP INDEX idx_rate_history_unique;