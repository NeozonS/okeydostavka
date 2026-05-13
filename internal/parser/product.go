package parser

import (
	"TestCenozavr/internal/client"
	"TestCenozavr/internal/models"
	"TestCenozavr/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"
)

func (p *Parser) GetProducts(category models.Category) (models.CategoryProducts, error) {
	var allProducts []models.Product
	pageSize := 72
	pageNumber := 1
	totalCount := 1

	for {
		select {
		case <-p.ctx.Done():
			return models.CategoryProducts{}, p.ctx.Err()
		default:
		}

		products, err := p.client.GetProductInfo(category.ID, pageNumber, pageSize)
		if err != nil {
			slog.Error("Error receiving products", "error", err, "category_id", category.ID, "page", pageNumber)
			return models.CategoryProducts{}, fmt.Errorf("failed to get products: %w", err)
		}

		var page models.Products

		if err := json.Unmarshal(products, &page); err != nil {
			slog.Error("Product parsing error", "error", err, "category_id", category.ID, "page", pageNumber, "body_length", len(products))
			return models.CategoryProducts{}, fmt.Errorf("failed to unmarshal products: %w", err)
		}

		if pageNumber == 1 {
			totalCount = page.TotalCount
		}
		allProducts = append(allProducts, page.Products...)

		if len(page.Products) < pageSize || len(allProducts) >= totalCount {
			break
		}
		pageNumber++

		utils.Jitter(800*time.Millisecond, 2500*time.Millisecond)
	}

	fullUrlProduct(allProducts)

	return models.CategoryProducts{
		CategoryID:   category.ID,
		CategoryName: category.Name,
		Products:     allProducts,
	}, nil
}

func (p *Parser) GetAllProducts() ([]models.CategoryProducts, error) {
	var allCategoriesProducts []models.CategoryProducts
	blockCount := 0

	cat, err := p.GetCatalog()
	if err != nil {
		slog.Error("Failed to get catalog", "error", err)
		return nil, fmt.Errorf("failed to get catalog: %w", err)
	}

	if len(cat.Categories) == 0 {
		slog.Warn("Catalog is empty, no categories to process")
		return allCategoriesProducts, nil
	}

	slog.Info("Found categories to process", "count", len(cat.Categories))

	rand.Shuffle(len(cat.Categories), func(i, j int) {
		cat.Categories[i], cat.Categories[j] = cat.Categories[j], cat.Categories[i]
	})

	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := p.client.GetCookie(); err != nil {
					slog.Error("Failed to get cookies", "error", err)
				}
			case <-p.ctx.Done():
				return
			}
		}
	}()

	for i, category := range cat.Categories {
		select {
		case <-p.ctx.Done():
			slog.Info("Parsing stopped", "processed_categories", len(allCategoriesProducts), "total_products", countProducts(allCategoriesProducts))
			return allCategoriesProducts, p.ctx.Err()
		default:
		}

		utils.Jitter(1500*time.Millisecond, 5*time.Second)

		prod, err := p.GetProducts(category)
		if err != nil {
			if errors.Is(err, client.AntiBotError) {
				blockCount++
				if blockCount >= 3 {
					slog.Warn("Blocked by anti-bot 3 times", "processed_categories", len(allCategoriesProducts), "total_products", countProducts(allCategoriesProducts))
					return allCategoriesProducts, fmt.Errorf("blocked 3 times: %w", client.AntiBotError)
				}
			} else {
				blockCount = 0
			}

			slog.Warn("Error in product search, skipping category", "category_id", category.ID, "category_name", category.Name, "error", err)
			continue
		}

		blockCount = 0

		if len(prod.Products) > 0 {
			allCategoriesProducts = append(allCategoriesProducts, prod)
			slog.Debug("Category processed", "current", i+1, "total", len(cat.Categories), "category_name", category.Name, "products", len(prod.Products))
		}
	}

	slog.Info("Parsing completed", "processed_categories", len(allCategoriesProducts), "total_products", countProducts(allCategoriesProducts))
	return allCategoriesProducts, nil
}

func countProducts(categories []models.CategoryProducts) int {
	total := 0
	for _, c := range categories {
		total += len(c.Products)
	}
	return total
}

func fullUrlProduct(prod []models.Product) {
	for i, _ := range prod {
		prod[i].ProductURL = "https://www.okeydostavka.ru" + prod[i].ProductURL
	}
}
