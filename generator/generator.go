package generator

import (
	"reflect"

	"github.com/go-faker/faker/v4"
)

type Generator struct {
	model  any
}

func NewGenerator(model any) *Generator {
	return &Generator{model: model}
}

func (g *Generator) Generate(amount int) (data []any, err error) {
	// generate data
	data = make([]any, amount)
	for i := 0; i < amount; i++ {
		m := reflect.New(reflect.TypeOf(g.model)).Elem().Interface()
		err = faker.FakeData(&m)
		if err != nil {
			return nil, err
		}
		data[i] = m
	}

	for i := 0; i < amount; i++ {
		m := reflect.New(reflect.TypeOf(g.model)).Elem().Interface()
		err = faker.FakeData(&m)
		if err != nil {
			return nil, err
		}
		data[i] = m
	}

	return data, nil
}

// func (g *Generator) GenerateReplace(amount int, value any) error {
// 	// generate data

// 	data, ok := value.(*[]reflect.Type)
// 	if !ok {
// 		return fmt.Errorf("value is not a pointer to slice")
// 	}

// 	*data = make([]any, amount)
// 	for i := 0; i < amount; i++ {
// 		m := reflect.New(reflect.TypeOf(g.model)).Elem().Interface()
// 		err := faker.FakeData(&m)
// 		if err != nil {
// 			return err
// 		}
// 		(*data)[i] = m
// 	}

// 	return nil
// }

