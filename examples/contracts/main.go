// Команда contracts демонстрирует создание договора через билдер.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/avito-tech/avito-ads-sdk-go/avito"
)

func main() {
	client, err := avito.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	builder := avito.IntermediaryContract().
		Advertiser(987654321).
		CounterpartyType(avito.CounterpartyDirectWithAdvertiser).
		Subject(avito.ContractSubjectMediation).
		Object(avito.ContractActionCommercial).
		ReportingRequired(true).
		FundsAllocationToPrincipal(false).
		Date("2025-01-15").
		Number("ДА-2025/01").
		Intermediary(map[string]any{
			"shortName":     "ООО Реклама",
			"longName":      "Общество с ограниченной ответственностью «Реклама»",
			"inn":           "7712345678",
			"ogrn":          "1177746123456",
			"kpp":           "771701001",
			"legalAddress":  "г. Москва, ул. Примерная, д. 1",
			"actualAddress": "г. Москва, ул. Примерная, д. 1",
			"legalType":     avito.LegalTypeUL,
		})

	created, err := client.Contracts.Create(ctx, builder)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Создан договор ID: %d\n", created.ID)
}
