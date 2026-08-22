package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetJSON — Verilen URL'e GET isteği atar ve cevabı target'a unmarshal eder
// Tüm client'lar bu yardımcı fonksiyonu kullanabilir
func GetJSON(ctx context.Context, client *http.Client, url string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("request oluşturulamadı: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("istek başarısız: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API hata kodu %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("JSON parse hatası: %w", err)
	}

	return nil
}
