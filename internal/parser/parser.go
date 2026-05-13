package parser

import (
	"TestCenozavr/internal/client"
	"TestCenozavr/internal/config"
	"context"
)

type Parser struct {
	client *client.Client
	config *config.Config
	ctx    context.Context
}

func NewParser(ctx context.Context, client *client.Client, config *config.Config) *Parser {
	return &Parser{
		client: client,
		config: config,
		ctx:    ctx,
	}
}
