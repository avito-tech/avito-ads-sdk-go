// Команда pagination демонстрирует постраничную выборку и обход всех страниц.
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

	// Постранично вручную.
	page, err := client.Campaigns.List(ctx, avito.ListRequest{Limit: 50, Page: 1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Всего кампаний: %d, остаток лимита API: %d\n", page.Total, page.APIPointBalance)
	for _, c := range page.Items {
		fmt.Printf("  - %d %s [%s]\n", c.ID, c.Name, c.Status)
	}

	// Автоматический обход всех страниц.
	count := 0
	err = client.Campaigns.Iterate(ctx, avito.ListRequest{}, func(c avito.Campaign) bool {
		count++
		return true // вернуть false, чтобы остановить обход досрочно
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Обойдено кампаний: %d\n", count)
}
