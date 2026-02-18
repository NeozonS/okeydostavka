package parser

import (
	"TestCenozavr/internal/client"
	"TestCenozavr/internal/config"
)

type Parser struct {
	client *client.ClientWithCookies
	config *config.Config
}

func NewParser(client *client.ClientWithCookies, config *config.Config) *Parser {
	return &Parser{
		client: client,
		config: config,
	}
}
