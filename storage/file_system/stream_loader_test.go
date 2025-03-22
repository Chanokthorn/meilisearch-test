package file_system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// declared locally for modularity
type Product struct {
	ID          string  `json:"id" faker:"uuid_hyphenated"`
	Name        string  `json:"name" faker:"first_name"`
	Price       float64 `json:"price" faker:"amount"`
	Description string  `json:"description" faker:"paragraph len=800"`
	Score       float64 `json:"score" faker:"boundary_start=-2, boundary_end=2"`
}

func Test_streamLoader_Start(t *testing.T) {

	t.Run("running", func(t *testing.T) {
		sl := NewStreamLoader(WithReadBatchSize(30), WithTesting())
		sl.SetModel(Product{})
		if err := sl.SetReadFile("products.output"); err != nil {
			t.Fatal(err)
		}

		dataChan, errChan := sl.Start()

		var dataResult []Product

	Loop:
		for {
			select {
			case data, ok := <-dataChan:
				if !ok {
					break Loop
				}
				product, ok := data.(Product)
				assert.True(t, ok, "data should be of type Product")

				dataResult = append(dataResult, product)
			case err := <-errChan:
				assert.NoError(t, err)
				return
			}
		}

		assert.Len(t, dataResult, 100)

	})
}
