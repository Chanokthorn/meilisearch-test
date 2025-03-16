package generator

import (
	"ms-tester/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator_Generate(t *testing.T) {
	t.Run("Generate", func(t *testing.T) {
		length := 200

		gn := NewGenerator(model.Product{}, length)
		data, err := gn.Generate()
		assert.NoError(t, err)

		assert.Len(t, data, length)
		for _, d := range data {
			assert.NotZero(t, d)
		}
	})
}
