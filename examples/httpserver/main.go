// Команда httpserver показывает интеграцию клиента в веб-сервис на стандартной
// библиотеке. Клиент создаётся один раз при старте и переиспользуется во всех
// обработчиках (он потокобезопасен). Тот же приём применим к любому фреймворку
// (Gin, Echo, Fiber, Chi) — см. README.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/avito-tech/avito-ads-sdk-go/avito"
)

func main() {
	client, err := avito.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/balance", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()

		balance, err := client.Account.GetBalance(ctx)
		if err != nil {
			status := http.StatusBadGateway
			var apiErr *avito.APIError
			if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusTooManyRequests {
				status = http.StatusTooManyRequests
			}
			http.Error(w, err.Error(), status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"balance": balance.Balance,
			"bonus":   balance.BonusBalance,
		})
	})

	log.Println("слушаю :8080")
	srv := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
