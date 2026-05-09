package main

import (
	"TestCenozavr/internal/client"
	"TestCenozavr/internal/config"
	"TestCenozavr/internal/utils"
	"TestCenozavr/internal/parser"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		})))

	slog.Info("Application started")

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Config error", "error", err)
		os.Exit(1)
	}

	cln := client.NewClientWithCookies(cfg)
	defer cln.Close()

	parser := parser.NewParser(cln, cfg)

	if err := parser.InitializeStoreID(); err != nil {
		slog.Error("Failed to initialize store ID", "error", err)
		os.Exit(1)
	}

	val, err := parser.GetAllProducts()
	if err != nil {
		slog.Error("Failed to get all products", "error", err)
		os.Exit(1)
	}
	if len(val) == 0 {
		slog.Error("No products collected, check logs for details")
		os.Exit(1)
	}
	err = utils.SaveResult(val, "AllCat.txt")
	if err != nil {
		slog.Error("Failed to save results", "error", err)
		os.Exit(1)
	}
}
	