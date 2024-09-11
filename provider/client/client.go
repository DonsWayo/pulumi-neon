package client

import (
	"context"
	"log"

	"github.com/kislerdm/neon-sdk-go"
)

type Client struct {
	sdk *neon.SDK
}

func NewClient(apiKey string) *Client {
	sdk, err := neon.NewSDK(neon.Config{
		Key: apiKey,
	})
	if err != nil {
		log.Fatalf("Failed to create Neon SDK client: %v", err)
	}

	return &Client{
		sdk: sdk,
	}
}

// GetContext returns a new context for SDK operations
func (c *Client) GetContext() context.Context {
	return context.Background()
}