package stroage

import (
	"context"
	"ms-tester/model"
	"os"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
)

func genProducts(n int) ([]model.Product, error) {
	products := make([]model.Product, n)
	for i := range products {
		err := faker.FakeData(&products[i])
		if err != nil {
			return nil, err
		}
	}
	return products, nil
}

func TestProductStorage_SaveProduct(t *testing.T) {
	products, err := genProducts(10)
	assert.NoError(t, err)

	ps := NewProductStorage()
	err = ps.SaveProduct(context.Background(), products, "products.data")
	assert.NoError(t, err)

	assert.NoError(t, os.Remove("products.data"))
}

func TestProductStorage_LoadProduct(t *testing.T) {
	products, err := genProducts(250)
	assert.NoError(t, err)

	ps := NewProductStorage()
	err = ps.SaveProduct(context.Background(), products, "products-read.data")
	assert.NoError(t, err)

	assert.NoError(t, ps.SetReadFile("products-read.data"))

	loaded1, next1, err := ps.LoadProduct(t.Context())
	assert.NoError(t, err)
	assert.Len(t, loaded1, 100)
	assert.True(t, next1)

	loaded2, next2, err := ps.LoadProduct(t.Context())
	assert.NoError(t, err)
	assert.Len(t, loaded2, 100)
	assert.True(t, next2)

	loaded3, next3, err := ps.LoadProduct(t.Context())
	assert.NoError(t, err)
	assert.Len(t, loaded3, 50)
	assert.False(t, next3)

	os.Remove("products-read.data")
}
