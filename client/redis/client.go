package redis

import (
	"github.com/RediSearch/redisearch-go/redisearch"
	"weather-api/config"
)

type Client struct {
	rsc *redisearch.Client
}

func NewClient(cfg config.Redis, indexName string) *Client {
	return &Client{
		rsc: redisearch.NewClient(
			cfg.Address,
			indexName),
	}
}

func (c *Client) Index(docs ...redisearch.Document) error {
	return c.rsc.Index(docs...)
}

func (c *Client) Search(q *redisearch.Query) ([]redisearch.Document, error) {
	docs, _, err := c.rsc.Search(q)
	return docs, err
}

func (c *Client) Delete(id string) error {
	return c.rsc.DeleteDocument(id)
}
