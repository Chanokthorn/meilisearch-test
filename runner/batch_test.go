package runner

import (
	"context"
	mocks "ms-tester/mocks/meilisearch"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_batchRunner_Run(t *testing.T) {
	t.Run("Test Batch Run successfully", func(t *testing.T) {
		amount := 1000
		data := generateItems(amount)
		loader := &mockStreamLoader{data: data}
		uploadCalls := make([]bool, amount)

		ms := mocks.NewMeiliSearch(t)
		ms.EXPECT().AddOrUpdateDocument(mock.Anything, "indexUid", mock.Anything).RunAndReturn(func(ctx context.Context, indexUid string, dataBatch any) (int, error) {
			var latestTaskID int
			for _, data := range dataBatch.([]any) {
				it := data.(item)
				uploadCalls[it.ID] = true
				latestTaskID = it.ID
			}
			return latestTaskID, nil
		})

		w := NewBatchWorker(ms).SetBatchSize(15).SetIndexUid("indexUid")
		rn := NewRunner().SetWorker(w).SetWorkerAmount(3)

		lastTaskID, err := rn.Run(context.Background(), loader)
		assert.NoError(t, err)
		for i := 0; i < amount; i++ {
			assert.True(t, uploadCalls[i])
		}
		assert.Equal(t, amount-1, lastTaskID)
	})
}
