-- +goose Up
-- Заполняет базу историей курсов за последние 7 дней.

INSERT INTO rate_history (base_currency, quote_currency, value, date, source) VALUES
('USD', 'RUB', '82.07', NOW() - INTERVAL '7 days', 'seed'),
('USD', 'RUB', '81.97', NOW() - INTERVAL '6 days', 'seed'),
('USD', 'RUB', '82.00', NOW() - INTERVAL '5 days', 'seed'),
('USD', 'RUB', '82.34', NOW() - INTERVAL '4 days', 'seed'),
('USD', 'RUB', '82.21', NOW() - INTERVAL '3 days', 'seed'),
('USD', 'RUB', '82.38', NOW() - INTERVAL '2 days', 'seed'),
('USD', 'RUB', '82.43', NOW() - INTERVAL '1 day', 'seed'),

('EUR', 'RUB', '89.98', NOW() - INTERVAL '7 days', 'seed'),
('EUR', 'RUB', '90.23', NOW() - INTERVAL '6 days', 'seed'),
('EUR', 'RUB', '90.15', NOW() - INTERVAL '5 days', 'seed'),
('EUR', 'RUB', '90.52', NOW() - INTERVAL '4 days', 'seed'),
('EUR', 'RUB', '90.31', NOW() - INTERVAL '3 days', 'seed'),
('EUR', 'RUB', '90.70', NOW() - INTERVAL '2 days', 'seed'),
('EUR', 'RUB', '90.74', NOW() - INTERVAL '1 day', 'seed'),

('CNY', 'RUB', '11.55', NOW() - INTERVAL '3 days', 'seed'),
('CNY', 'RUB', '11.50', NOW() - INTERVAL '2 days', 'seed'),
('CNY', 'RUB', '11.40', NOW() - INTERVAL '1 day', 'seed')
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM rate_history WHERE source = 'seed';