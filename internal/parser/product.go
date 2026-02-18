package parser

import (
	"TestCenozavr/internal/models"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

func (p *Parser) GetProducts(category models.Category) (models.CategoryProducts, error) {
	var allProducts []models.Product
	pageSize := 72
	pageNumber := 1
	totalCount := 1

	for {
		products, err := p.client.GetProductInfo(category.ID, pageNumber, pageSize)
		if err != nil {
			slog.Error("Error receiving products", "error", err, "category_id", category.ID, "page", pageNumber)
			return models.CategoryProducts{}, fmt.Errorf("failed to get products: %w", err)
		}

		var page models.Products

		if err := json.Unmarshal(products.BodyBytes, &page); err != nil {
			slog.Error("Product parsing error", "error", err, "category_id", category.ID, "page", pageNumber, "body_length", len(products.BodyBytes))
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

		time.Sleep(1 * time.Second)
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

	if err := p.client.GetCookie(); err != nil {
		slog.Error("Failed to get cookies", "error", err)
	}

	cat, err := p.GetCatalog()
	if err != nil {
		slog.Error("Failed to get catalog", "error", err)
		return nil, fmt.Errorf("failed to get catalog: %w", err)
	}

	if len(cat.Categories) == 0 {
		slog.Warn("Catalog is empty, no categories to process")
		return allCategoriesProducts, nil
	}

	ch := make(chan struct{})
	defer close(ch)
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := p.client.GetCookie(); err != nil {
					slog.Error("Failed to get cookies", "error", err)
				}

			case <-ch:
				return
			}
		}
	}()

	for _, category := range cat.Categories {
		time.Sleep(1 * time.Second)

		prod, err := p.GetProducts(category)
		if err != nil {
			slog.Warn("Error in product search, skipping category", "category_id", category.ID, "category_name", category.Name, "error", err)
			continue
		}

		if len(prod.Products) > 0 {
			allCategoriesProducts = append(allCategoriesProducts, prod)
		}
	}

	return allCategoriesProducts, nil
}

func fullUrlProduct(prod []models.Product) {
	for i, _ := range prod {
		prod[i].ProductURL = "https://www.okeydostavka.ru" + prod[i].ProductURL
	}
}
