package storage

import (
	"context"
	"ms-tester/model"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_Save(t *testing.T) {
	products, err := genProducts(10)
	assert.NoError(t, err)

	anyArray := make([]any, len(products))
	for i, v := range products {
		anyArray[i] = v
	}

	st := NewStorage()
	err = st.Save(context.Background(), anyArray, "products.data")
	assert.NoError(t, err)

	assert.NoError(t, os.Remove("products.data"))
}

func TestStorage_Load(t *testing.T) {
	products, err := genProducts(250)
	assert.NoError(t, err)

	anyArray := make([]any, len(products))
	for i, v := range products {
		anyArray[i] = v
	}

	st := NewStorage().SetModel(model.Product{})
	err = st.Save(context.Background(), anyArray, "products-read.data")
	assert.NoError(t, err)

	assert.NoError(t, st.SetReadFile("products-read.data"))

	loaded1, next1, err := st.Load(t.Context())
	assert.NoError(t, err)
	assert.Len(t, loaded1, 100)
	assert.True(t, next1)

	loaded2, next2, err := st.Load(t.Context())
	assert.NoError(t, err)
	assert.Len(t, loaded2, 100)
	assert.True(t, next2)

	loaded3, next3, err := st.Load(t.Context())
	assert.NoError(t, err)
	assert.Len(t, loaded3, 50)
	assert.False(t, next3)

	os.Remove("products-read.data")
}
