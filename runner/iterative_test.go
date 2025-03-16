package runner

import (
	"ms-tester/mocks/meilisearch"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type item struct {
	ID          string  `json:"id"`
}

func TestIterative_Run(t *testing.T) {
	t.Run("Test Iterative Run", func(t *testing.T) {
		data := []any {
			item{
				ID: "1",
			},
			item{
				ID: "2",
			},
			item{
				ID: "3",
			},
			item{
				ID: "4",
			},
		}

		ms := mocks.NewMeiliSearch(t)
		ms.EXPECT().AddOrUpdateDocument(mock.Anything, mock.Anything, mock.Anything).Return(0, nil).Times(1)
		ms.EXPECT().AddOrUpdateDocument(mock.Anything, mock.Anything, mock.Anything).Return(1, nil).Times(1)
		ms.EXPECT().AddOrUpdateDocument(mock.Anything, mock.Anything, mock.Anything).Return(2, nil).Times(1)
		ms.EXPECT().AddOrUpdateDocument(mock.Anything, mock.Anything, mock.Anything).Return(3, nil).Times(1)
		it := NewIterative(ms).SetIndexUid("indexUid").SetData(data)
		lastID, err := it.Run(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, 3, lastID)
		
	})
}
