package meilisearch

import (
	"ms-tester/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMeiliSearch_CreateIndex(t *testing.T) {

	t.Run("Test CreateIndex", func(t *testing.T) {
		m := NewMeiliSearch("http://localhost:7700", "MASTER_KEY")
		err := m.CreateIndex(t.Context(), "test-index", "primaryKey")
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
