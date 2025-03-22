package runner

import (
	"context"
	mocks "ms-tester/mocks/meilisearch"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type item struct {
	ID int `json:"id"`
}

func generateItems(amount int) []item {
	items := make([]item, amount)
	for i := 0; i < amount; i++ {
		items[i] = item{
			ID: i,
		}
	}
	return items
}

type mockStreamLoader struct {
	data []item
}

func (m *mockStreamLoader) SetModel(model any) {}

func (m *mockStreamLoader) Start() (<-chan any, <-chan error) {
	dataChan := make(chan any, 10)
	errChan := make(chan error, 10)

	go func() {
		for _, item := range m.data {
			dataChan <- item
		}
		close(dataChan)
		close(errChan)
	}()

	return dataChan, errChan
}

func TestIterative_Run(t *testing.T) {
	t.Run("Test Iterative Run successfully", func(t *testing.T) {
		amount := 1000
		data := generateItems(amount)
		loader := &mockStreamLoader{data: data}
		uploadCalls := make([]bool, amount)

		ms := mocks.NewMeiliSearch(t)
		ms.EXPECT().AddOrUpdateDocument(mock.Anything, "indexUid", mock.Anything).RunAndReturn(func(ctx context.Context, indexUid string, data any) (int, error) {
			it := data.(item)
			uploadCalls[it.ID] = true
			return it.ID, nil
		})

		wk := NewIterativeWorker(ms).SetIndexUid("indexUid")
		rn := NewRunner().SetWorker(wk).SetWorkerAmount(4)

		ctx, _ := context.WithTimeout(t.Context(), 3*time.Second)

		lastID, err := rn.Run(ctx, loader)
		assert.NoError(t, err)
		assert.Equal(t, amount-1, lastID)
		for _, called := range uploadCalls {
			assert.True(t, called)
		}

	})
}
