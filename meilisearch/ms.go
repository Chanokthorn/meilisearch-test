package meilisearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"gopkg.in/resty.v1"
)

type MeiliSearch interface {
	CreateIndex(ctx context.Context, indexID string, primaryKey string) error
	AddOrUpdateDocument(ctx context.Context, indexName string, document any) (taskUid int, err error)
}

type meiliSearch struct {
	client *resty.Client
}

func NewMeiliSearch(endpoint string, apiKey string) *meiliSearch {
	client := resty.New()
	client.SetHostURL(endpoint)
	client.SetHeader("content-type", "application/json")
	client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	return &meiliSearch{
		client: client,
	}
}

func (m *meiliSearch) CreateIndex(ctx context.Context, indexID string, primaryKey string) error {
	_, err := m.client.R().SetBody(map[string]interface{}{
		"uid":        indexID,
		"primaryKey": primaryKey,
	}).Post("/indexes")

	if err != nil {
		err := fmt.Errorf("error creating index: %w", err)
		log.Printf("%s", err.Error())
		return err
	}

	return nil
}

type AddOrUpdateDocumentResponse struct {
	TaskUid int `json:"taskUid"`
	// Status string `json:"status"` // "enqueued"
	// ...
}

func (m *meiliSearch) AddOrUpdateDocument(ctx context.Context, indexName string, document any) (taskUid int, err error) {
	res, err := m.client.R().SetBody(document).Post(fmt.Sprintf("/indexes/%s/documents", indexName))
	if err != nil {
		err := fmt.Errorf("error adding or updating document: %w", err)
		log.Printf("%s", err.Error())
		return 0, err
	}

	if res.StatusCode() != 202 {
		err := fmt.Errorf("error adding or updating document: %s", res.String())
		log.Printf("%s", err.Error())
		return 0, err
	}

	// parse response
	var result AddOrUpdateDocumentResponse
	if err := json.Unmarshal(res.Body(), &result); err != nil {
		log.Printf("error parsing response body: %s", err.Error())
		return 0, err
	}

	return result.TaskUid, nil
}
