package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

type Storage struct {
	readBatchSize int
	file          *os.File
	scanner       *bufio.Scanner
	model         reflect.Type
}

func NewStorage() *Storage {
	return &Storage{
		readBatchSize: 100,
	}
}

func (st *Storage) SetModel(model any) *Storage {
	st.model = reflect.TypeOf(model)
	return st
}

func (st *Storage) Save(ctx context.Context, items []any, path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err := fmt.Errorf("error creating file: %w", err)
		log.Printf("%s", err.Error())
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, i := range items {
		data, err := json.Marshal(i)
		if err != nil {
			err := fmt.Errorf("error marshalling item: %w", err)
			log.Printf("%s", err.Error())
			return err
		}

		_, err = writer.Write(data)
		if err != nil {
			err := fmt.Errorf("error writing item: %w", err)
			log.Printf("%s", err.Error())
			return err
		}

		writer.WriteString("\n")
	}

	if err := writer.Flush(); err != nil {
		err := fmt.Errorf("error flushing writer: %w", err)
		log.Printf("%s", err.Error())
		return err
	}

	return nil
}

func (st *Storage) SetReadFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		err := fmt.Errorf("error opening file: %w", err)
		log.Printf("%s", err.Error())
		return err
	}

	st.file = file
	st.scanner = bufio.NewScanner(st.file)

	return nil
}

func (st *Storage) Load(ctx context.Context) ([]any, bool, error) {
	if st.file == nil || st.scanner == nil {
		err := fmt.Errorf("file or scanner not set")
		log.Printf("%s", err.Error())
		return nil, false, err
	}

	var items []any
	counter := 0
	var end bool
	for counter < st.readBatchSize {
		if !st.scanner.Scan() {
			end = true
			break
		}
		item := reflect.New(st.model).Interface()
		if err := json.Unmarshal(st.scanner.Bytes(), &item); err != nil {
			err := fmt.Errorf("error unmarshalling product: %w", err)
			log.Printf("%s", err.Error())
			return nil, false, err
		}

		items = append(items, item)
		counter++
	}

	if end {
		st.file.Close()
	}

	return items, !end, nil
}

// func (st *Storage) Start(ctx context.Context) (error) {
// 	if st.file == nil || st.scanner == nil {
// 		err := fmt.Errorf("file or scanner not set")
// 		log.Printf("%s", err.Error())
// 		return nil, false, err

// 	}

// 	counter := 0
// 	for counter < st.readBatchSize {
// 		if !st.scanner.Scan() {
// 			break
// 		}
// 		item := reflect.New(st.model).Interface()
// 		if err := json.Unmarshal(st.scanner.Bytes(), item); err != nil {
// 			err := fmt.Errorf("error unmarshalling item: %w", err)
// 			log.Printf("%s", err.Error())
// 			return nil, false, err
// 		}

// 		items = append(items, product)
// 		counter++
// 	}

// 	st.file.Close()

// 	return items, !end, nil
// }
