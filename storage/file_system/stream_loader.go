package file_system

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
)

type streamLoader struct {
	file          *os.File
	scanner       *bufio.Scanner
	readBatchSize int
	model         reflect.Type
	testing       bool
}

type streamLoaderOption func(sl streamLoader) streamLoader

func WithTesting() streamLoaderOption {
	return func(sl streamLoader) streamLoader {
		sl.testing = true
		return sl
	}
}

func WithReadBatchSize(size int) streamLoaderOption {
	return func(sl streamLoader) streamLoader {
		sl.readBatchSize = size
		return sl
	}
}

func NewStreamLoader(options ...streamLoaderOption) *streamLoader {
	sl := streamLoader{
		readBatchSize: 100,
	}
	for _, option := range options {
		sl = option(sl)
	}
	return &sl
}

func getAbsPath(path string) (string, error) {
	return filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), path))
}

func (sl *streamLoader) SetReadFile(path string) error {
	var actualPath string
	if sl.testing {
		actualPath = path
	} else {
		fullPath, err := getAbsPath(path)
		if err != nil {
			err := fmt.Errorf("error getting absolute path: %w", err)
			log.Printf("%s", err.Error())
			return err
		}

		actualPath = fullPath
	}

	file, err := os.OpenFile(actualPath, os.O_RDONLY, 0644)
	if err != nil {
		err := fmt.Errorf("error opening file: %w", err)
		log.Printf("%s", err.Error())
		return err
	}

	sl.file = file

	sl.scanner = bufio.NewScanner(sl.file)

	return nil
}

func (sl *streamLoader) SetModel(model any) {
	sl.model = reflect.TypeOf(model)
}

func (sl *streamLoader) Start() (<-chan any, <-chan error) {
	if sl.model == nil {
		log.Fatalf("model not set")
	}

	ctx := context.Background()
	out := make(chan any)
	errChan := make(chan error)

	go func() {
		for {
			dataBatch, next, err := sl.load(ctx)
			if err != nil {
				log.Printf("error loading data: %s", err.Error())
				errChan <- fmt.Errorf("error loading data: %w", err)
				return
			}

			for _, data := range dataBatch {
				out <- reflect.ValueOf(data).Elem().Interface()
			}

			if !next {
				break
			}
		}

		close(out)
		close(errChan)
	}()

	return out, errChan
}

func (sl *streamLoader) load(ctx context.Context) ([]any, bool, error) {
	var items []any
	counter := 0
	var end bool
	for counter < sl.readBatchSize {
		if !sl.scanner.Scan() {
			end = true
			break
		}
		item := reflect.New(sl.model).Interface()
		if err := json.Unmarshal(sl.scanner.Bytes(), &item); err != nil {
			err := fmt.Errorf("error unmarshalling product: %w", err)
			log.Printf("%s", err.Error())
			return nil, false, err
		}

		items = append(items, item)
		counter++
	}

	if end {
		err := sl.file.Close()
		if err != nil {
			err := fmt.Errorf("error closing file: %w", err)
			log.Printf("%s", err.Error())
			return nil, false, err
		}
	}

	return items, !end, nil
}
