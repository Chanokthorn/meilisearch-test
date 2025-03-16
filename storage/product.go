package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ms-tester/model"
	"os"
)

type ProductStorage struct {
	readBatchSize int
	file          *os.File
	scanner       *bufio.Scanner
}

func NewProductStorage() *ProductStorage {
	return &ProductStorage{
		readBatchSize: 100,
	}
}

func (p *ProductStorage) SaveProduct(ctx context.Context, products []model.Product, path string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err := fmt.Errorf("error creating file: %w", err)
		log.Printf("%s", err.Error())
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, p := range products {
		data, err := json.Marshal(p)
		if err != nil {
			err := fmt.Errorf("error marshalling product: %w", err)
			log.Printf("%s", err.Error())
			return err
		}

		_, err = writer.Write(data)
		if err != nil {
			err := fmt.Errorf("error writing product: %w", err)
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

func (p *ProductStorage) SetReadFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		err := fmt.Errorf("error opening file: %w", err)
		log.Printf("%s", err.Error())
		return err
	}

	p.file = file
	p.scanner = bufio.NewScanner(p.file)

	return nil
}

func (p *ProductStorage) LoadProduct(ctx context.Context) ([]model.Product, bool, error) {
	if p.file == nil || p.scanner == nil {
		err := fmt.Errorf("file or scanner not set")
		log.Printf("%s", err.Error())
		return nil, false, err

	}

	products := make([]model.Product, 0, p.readBatchSize)
	counter := 0
	var end bool
	for counter < p.readBatchSize {
		if !p.scanner.Scan() {
			end = true
			break
		}
		var product model.Product
		if err := json.Unmarshal(p.scanner.Bytes(), &product); err != nil {
			err := fmt.Errorf("error unmarshalling product: %w", err)
			log.Printf("%s", err.Error())
			return nil, false, err
		}

		products = append(products, product)
		counter++
	}

	if end {
		p.file.Close()
	}

	return products, !end, nil
}
