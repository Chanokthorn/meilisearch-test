package meilisearch

import (
	"ms-tester/model"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

func TestMeiliSearch_CreateIndex(t *testing.T) {

	t.Run("Test CreateIndex", func(t *testing.T) {
		m := NewMeiliSearch("http://localhost:7700", "MASTER_KEY")
		err := m.CreateIndex(t.Context(), "test-index", "primaryKey")
		assert.NoError(t, err)
	})
}

func TestMeiliSearch_DeleteIndex(t *testing.T) {
	t.Run("Test DeleteIndex", func(t *testing.T) {
		m := NewMeiliSearch("http://localhost:7700", "MASTER_KEY")
		err := m.CreateIndex(t.Context(), "test-index-to-delete", "primaryKey")
		assert.NoError(t, err)

		err = m.DeleteIndex(t.Context(), "test-index-to-delete")
		assert.NoError(t, err)
	})
}

func TestMeiliSearch_AddOrUpdateDocument(t *testing.T) {

	indexName := "test-product"

	t.Run("Test AddOrUpdateDocument", func(t *testing.T) {
		document := model.Product{
			ID:          "23124",
			Name:        "john",
			Price:       123123,
			Description: "describeeeee",
			Score:       1.283210,
		}
		m := NewMeiliSearch("http://localhost:7700", "MASTER_KEY")

		err := m.CreateIndex(t.Context(), indexName, "id")
		assert.NoError(t, err)

		taskUid, err := m.AddOrUpdateDocument(t.Context(), indexName, document)
		assert.NoError(t, err)
		assert.NotZero(t, taskUid)
	})

	t.Run("Test multiple documents", func(t *testing.T) {
		documents := []model.Product{
			{
				ID:          "231",
				Name:        "john2",
				Price:       123123,
				Description: "describeeeee",
				Score:       -1.283210,
			},
			{
				ID:          "2314",
				Name:        "john4",
				Price:       123123,
				Description: "describeeeee",
				Score:       1.183210,
			},
			{
				ID:          "2312423",
				Name:        "john5",
				Price:       123123,
				Description: "describeeeee",
				Score:       1.21210,
			},
		}
		m := NewMeiliSearch("http://localhost:7700", "MASTER_KEY")

		err := m.CreateIndex(t.Context(), indexName, "id")
		assert.NoError(t, err)

		taskUid, err := m.AddOrUpdateDocument(t.Context(), indexName, documents)
		assert.NoError(t, err)
		assert.NotZero(t, taskUid)
	})
}

func Test_meiliSearch_getTaskStatus(t *testing.T) {
	t.Run("Test getTaskStatus", func(t *testing.T) {
		m := NewMeiliSearch("http://localhost:7700", "MASTER_KEY")

		taskUid, err := m.AddOrUpdateDocument(t.Context(), "test-product", model.Product{
			ID:          "23124",
			Name:        "john",
			Price:       123123,
			Description: "describeeeee",
			Score:       1.283210,
		})
		assert.NoError(t, err)

		status1, err := m.getTaskStatus(t.Context(), taskUid)
		assert.NoError(t, err)
		assert.Contains(t, []string{"enqueued", "processing"}, status1)

		time.Sleep(500 * time.Millisecond)

		status2, err := m.getTaskStatus(t.Context(), taskUid)
		assert.NoError(t, err)
		assert.Equal(t, "succeeded", status2)

	})
}

type item struct {
	ID   string `json:"id" faker:"uuid_hyphenated"`
	Text string `json:"text" faker:"paragraph len=3000"`
}

func Test_meiliSearch_WaitTaskDone(t *testing.T) {
	items := make([]item, 100)
	for i := range items {
		faker.FakeData(&items[i])
	}

	t.Run("Test WaitTaskDone", func(t *testing.T) {
		m := NewMeiliSearch("http://localhost:7700", "MASTER_KEY")

		err := m.CreateIndex(t.Context(), "test-index-wait-task-done", "id")
		assert.NoError(t, err)

		taskUid, err := m.AddOrUpdateDocument(t.Context(), "test-index-wait-task-done", items)
		assert.NoError(t, err)

		t.Log("enqueued...")

		err = m.WaitTaskDone(t.Context(), taskUid)
		assert.NoError(t, err)
	})

}
