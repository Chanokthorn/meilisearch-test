package main

import (
	"context"
	"fmt"
	"ms-tester/cmd/cli/config"
	"ms-tester/generator"
	"ms-tester/meilisearch"
	"ms-tester/model"
	"ms-tester/runner"
)

func main() {
	cfg := config.Read()

	ms := meilisearch.NewMeiliSearch(cfg.Host, cfg.MasterKey)

	generator := generator.NewGenerator(model.Product{})

	data, err := generator.Generate(200)
	if err != nil {
		panic(fmt.Errorf("failed to generate data: %w", err))
	}

	ctx := context.Background()

	ms.CreateIndex(ctx, "products", "id")

	runner := runner.NewIterative(ms).SetIndexUid("products").SetData(data)

	lastID, err := runner.Run(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to run: %w", err))
	}

	println("Last ID:", lastID)

	// generate data

	// start report generator

	// start runner with mode iterative or batch

	// gen report
}