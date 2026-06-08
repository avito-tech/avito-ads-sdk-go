// Команда quickstart демонстрирует получение баланса и реквизитов аккаунта.
//
// Запуск:
//
//	AVITO_ADS_CLIENT_ID=... AVITO_ADS_CLIENT_SECRET=... AVITO_ADS_ACCOUNT_ID=... \
//	AVITO_ADS_ENVIRONMENT=sandbox go run ./examples/quickstart
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/avito-tech/avito-ads-sdk-go/avito"
)

func main() {
	client, err := avito.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	account, err := client.Account.Get(ctx)
	if err != nil {
		var apiErr *avito.APIError
		if errors.As(err, &apiErr) {
			log.Fatalf("ошибка API %d: %s", apiErr.StatusCode, apiErr.Message)
		}
		log.Fatal(err)
	}
	fmt.Printf("Аккаунт: %s (ИНН %s)\n", account.ShortName, account.INN)

	balance, err := client.Account.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Баланс: %d ₽, бонусы: %d\n", balance.Balance, balance.BonusBalance)
}
