package file

//nolint:depguard
import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/ArtemQa8/currency-converter/internal/model"
)

//go:embed seed_rates.json
var seedData []byte

// Repository описывает файловый репозиторий.
type Repository struct {
	filePath string
	mu       sync.RWMutex
}

// NewRepository конструктор - создаёт новый файловый репозиторий.
func NewRepository(filePath string) *Repository {
	return &Repository{filePath: filePath}
}

// SaveRate сохраняет курс в файл.
func (r *Repository) SaveRate(ctx context.Context, rate model.Rate) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data := make(map[string]model.Rate)
	fileBytes, err := os.ReadFile(r.filePath)
	if err == nil {
		_ = json.Unmarshal(fileBytes, &data)
	}

	key := fmt.Sprintf("%s_%s", rate.BaseCurrency, rate.QuoteCurrency)
	data[key] = rate
	newBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка маршалинга: %w", err)
	}

	tmpPath := r.filePath + ".tmp"
	tmpFile, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("ошибка открытия tmp файла: %w", err)
	}

	if _, err := tmpFile.Write(newBytes); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("ошибка записи в tmp файл: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("ошибка синхроназции: %w", err)
	}

	_ = tmpFile.Close()

	if err := os.Rename(tmpPath, r.filePath); err != nil {
		return fmt.Errorf("ошибка переименования: %w", err)
	}

	return nil
}

// GetRate получает курс из файла (или из seed, если файла нет).
func (r *Repository) GetRate(ctx context.Context, baseCurrency, quoteCurrency string) (model.Rate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var data map[string]model.Rate
	key := fmt.Sprintf("%s_%s", baseCurrency, quoteCurrency)

	fileBytes, err := os.ReadFile(r.filePath)
	if err == nil {
		if err := json.Unmarshal(fileBytes, &data); err == nil {
			if rate, ok := data[key]; ok {
				rate.Source = "file"
				return rate, nil
			}
		}
	}

	if err := json.Unmarshal(seedData, &data); err == nil {
		if rate, ok := data[key]; ok {
			return rate, nil
		}
	}

	return model.Rate{}, fmt.Errorf("курс %s не найден в файле", key)
}
