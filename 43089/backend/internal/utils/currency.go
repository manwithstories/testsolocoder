package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"travel-planner/internal/logger"
)

type ExchangeRateService struct {
	rates        map[string]float64
	lastUpdated  time.Time
	mu           sync.RWMutex
	baseCurrency string
}

var (
	exchangeRateService *ExchangeRateService
	once                sync.Once
)

const (
	exchangeRateAPI = "https://api.exchangerate-api.com/v4/latest/%s"
	cacheDuration   = 1 * time.Hour
)

func GetExchangeRateService() *ExchangeRateService {
	once.Do(func() {
		exchangeRateService = &ExchangeRateService{
			rates:        make(map[string]float64),
			baseCurrency: "CNY",
		}
		exchangeRateService.rates["CNY"] = 1.0
		exchangeRateService.rates["USD"] = 0.1389
		exchangeRateService.rates["EUR"] = 0.1285
		exchangeRateService.rates["JPY"] = 21.5
		exchangeRateService.rates["GBP"] = 0.1085
		exchangeRateService.rates["HKD"] = 1.08
		exchangeRateService.rates["TWD"] = 4.45
		exchangeRateService.rates["KRW"] = 195.0
		exchangeRateService.rates["THB"] = 4.95
		exchangeRateService.rates["SGD"] = 0.186
		exchangeRateService.rates["MYR"] = 0.66
		exchangeRateService.lastUpdated = time.Now()
	})
	return exchangeRateService
}

func (s *ExchangeRateService) fetchRates() error {
	url := fmt.Sprintf(exchangeRateAPI, s.baseCurrency)
	resp, err := http.Get(url)
	if err != nil {
		logger.Warnf("Failed to fetch exchange rates: %v, using fallback rates", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result struct {
		Rates map[string]float64 `json:"rates"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.rates = result.Rates
	s.lastUpdated = time.Now()
	logger.Infof("Exchange rates updated successfully")
	return nil
}

func (s *ExchangeRateService) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return amount, nil
	}

	if time.Since(s.lastUpdated) > cacheDuration {
		go s.fetchRates()
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	fromRate, ok := s.rates[fromCurrency]
	if !ok {
		return 0, fmt.Errorf("unsupported currency: %s", fromCurrency)
	}

	toRate, ok := s.rates[toCurrency]
	if !ok {
		return 0, fmt.Errorf("unsupported target currency: %s", toCurrency)
	}

	amountInBase := amount / fromRate
	converted := amountInBase * toRate

	return converted, nil
}

func (s *ExchangeRateService) ConvertToBase(amount float64, fromCurrency string) (float64, error) {
	return s.Convert(amount, fromCurrency, s.baseCurrency)
}

func (s *ExchangeRateService) GetSupportedCurrencies() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	currencies := make([]string, 0, len(s.rates))
	for currency := range s.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func ConvertCurrency(amount float64, fromCurrency, toCurrency string) (float64, error) {
	service := GetExchangeRateService()
	return service.Convert(amount, fromCurrency, toCurrency)
}

func ConvertToBaseCurrency(amount float64, fromCurrency string) (float64, error) {
	service := GetExchangeRateService()
	return service.ConvertToBase(amount, fromCurrency)
}
