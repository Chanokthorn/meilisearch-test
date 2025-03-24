package meilisearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gopkg.in/resty.v1"
)

type MeiliSearch interface {
	CreateIndex(ctx context.Context, indexID string, primaryKey string) error
	DeleteIndex(ctx context.Context, indexID string) error
	AddOrUpdateDocument(ctx context.Context, indexName string, document any) (taskUid int, err error)
	WaitTaskDone(ctx context.Context, taskUid int) error
}

type meiliSearch struct {
	client               *resty.Client
	statusUpdateInterval time.Duration
}

func NewMeiliSearch(endpoint string, apiKey string) *meiliSearch {
	client := resty.New()
	client.SetHostURL(endpoint)
	client.SetHeader("content-type", "application/json")
	client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	return &meiliSearch{
		client:               client,
		statusUpdateInterval: 100 * time.Millisecond,
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

func (m *meiliSearch) DeleteIndex(ctx context.Context, indexID string) error {
	_, err := m.client.R().Delete(fmt.Sprintf("/indexes/%s", indexID))
	if err != nil {
		err := fmt.Errorf("error deleting index: %w", err)
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
	dataBytes, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		log.Printf("error marshaling document to json: %s", err.Error())
		return 0, err
	}

	log.Printf("document: %s", string(dataBytes))
	
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

type GetTaskResponse struct {
	Status string `json:"status"`
	// ...
}

func (m *meiliSearch) getTaskStatus(ctx context.Context, taskUid int) (string, error) {
	res, err := m.client.R().Get(fmt.Sprintf("/tasks/%d", taskUid))
	if err != nil {
		err := fmt.Errorf("error getting task status: %w", err)
		log.Printf("%s", err.Error())
		return "", err
	}

	if res.StatusCode() != 200 {
		err := fmt.Errorf("error getting task status: %s", res.String())
		log.Printf("%s", err.Error())
		return "", err
	}

	// parse response
	var result GetTaskResponse
	if err := json.Unmarshal(res.Body(), &result); err != nil {
		log.Printf("error parsing response body: %s", err.Error())
		return "", err
	}

	return result.Status, nil
}

func (m *meiliSearch) WaitTaskDone(ctx context.Context, taskUid int) error {
	done := make(chan struct{})
	go func() {
		defer close(done)

		for {
			status, err := m.getTaskStatus(ctx, taskUid)
			if err != nil {
				log.Printf("error getting task status: %s", err.Error())
				return
			}

			if status == "succeeded" {
				return
			}

			time.Sleep(m.statusUpdateInterval)
		}
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
