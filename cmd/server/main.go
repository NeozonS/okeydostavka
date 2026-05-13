package main

import (
	"TestCenozavr/internal/client"
	"TestCenozavr/internal/config"
	"TestCenozavr/internal/parser"
	"TestCenozavr/internal/utils"
	"context"
	"errors"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	listStores := flag.Bool("list-stores", false, "List all available stores and exit")
	flag.Parse()

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

	cln, err := client.NewClientWithCookies(cfg)
	if err != nil {
		slog.Error("Failed to initialize client", "error", err)
		os.Exit(1)
	}
	defer cln.Close()

	if *listStores {
		ctx := context.Background()
		p := parser.NewParser(ctx, cln, cfg)

		if err := p.PrintAllStores(); err != nil {
			slog.Error("Failed to get stores", "error", err)
			os.Exit(1)
		}
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		slog.Info("Received shutdown signal", "signal", sig.String())
		cancel()
	}()

	p := parser.NewParser(ctx, cln, cfg)

	if err := p.InitializeStoreID(); err != nil {
		slog.Error("Failed to initialize store ID", "error", err)
		os.Exit(1)
	}

	val, err := p.GetAllProducts()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			slog.Info("Parsing interrupted by user")
		} else {
			slog.Error("Failed to get all products", "error", err)
		}
	}

	if len(val) > 0 {
		totalProducts := 0
		for _, cp := range val {
			totalProducts += len(cp.Products)
		}
		if err := utils.SaveResult(val, "AllCat.txt"); err != nil {
			slog.Error("Failed to save results", "error", err)
			os.Exit(1)
		}
		slog.Info("Results saved", "categories", len(val), "products", totalProducts)
	} else if err == nil {
		slog.Error("No products collected, check logs for details")
		os.Exit(1)
	}

	if err != nil && !errors.Is(err, context.Canceled) {
		os.Exit(1)
	}
}
