package generator

import (
	"reflect"

	"github.com/go-faker/faker/v4"
)

type Generator struct {
	model  any
	amount int
}

func NewGenerator(model any, amount int) *Generator {
	return &Generator{model: model, amount: amount}
}

func (g *Generator) Generate() (data []any, err error) {
	// generate data
	data = make([]any, g.amount)
	for i := 0; i < g.amount; i++ {
		m := reflect.New(reflect.TypeOf(g.model)).Elem().Interface()
		err = faker.FakeData(&m)
		if err != nil {
			return nil, err
		}
		data[i] = m
	}

	for i := 0; i < g.amount; i++ {
		m := reflect.New(reflect.TypeOf(g.model)).Elem().Interface()
		err = faker.FakeData(&m)
		if err != nil {
			return nil, err
		}
		data[i] = m
	}

	return data, nil
}
