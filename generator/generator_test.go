package generator

import (
	"ms-tester/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator_Generate(t *testing.T) {
	t.Run("Generate", func(t *testing.T) {
		length := 200
		gn := NewGenerator(model.Product{})
		data, err := gn.Generate(200)
		assert.NoError(t, err)

		assert.Len(t, data, length)
		for _, d := range data {
			assert.NotZero(t, d)
		}
	})
}

// func TestGenerator_GenerateReplace(t *testing.T) {
// 	t.Run("Generate", func(t *testing.T) {
// 		length := 200
// 		gn := NewGenerator(model.Product{})

// 		var data []model.Product
// 		err := gn.GenerateReplace(length, &data)
// 		assert.NoError(t, err)

// 		assert.Len(t, data, length)
// 		for _, d := range data {
// 			assert.NotZero(t, d)
// 		}
// 	})
// }
