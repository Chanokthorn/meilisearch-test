package pg

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type streamLoader struct {
	db          *sqlx.DB
	table       string
	primaryKey  string
	currOffset  int
	queryLimit  int
	sampleLimit int
}

func NewStreamLoader(url string) (*streamLoader, error) {
	pgConn, err := sqlx.Connect("postgres", url)
	if err != nil {
		err := fmt.Errorf("unable to connect to postgres: %w", err)
		return nil, err
	}

	return &streamLoader{
		db:         pgConn,
		queryLimit: 100,
		sampleLimit: 1000,
	}, nil
}

func (sl *streamLoader) SetTable(table string, primaryKey string) *streamLoader {
	sl.table = table
	sl.primaryKey = primaryKey
	return sl
}

func (sl *streamLoader) SetQueryLimit(limit int) *streamLoader {
	sl.queryLimit = limit
	return sl
}

func (sl *streamLoader) SetSampleLimit(limit int) *streamLoader {
	sl.sampleLimit = limit
	return sl
}

func (sl *streamLoader) Start() (chan interface{}, chan error) {
	dataChan := make(chan interface{}, 10)
	errChan := make(chan error, 10)

	ctx := context.Background()
	go func() {
		for {
			data, next, err := sl.load(ctx)
			if err != nil {
				errChan <- err
			}
			for _, data := range data {
				dataChan <- data
			}
			if !next {
				break
			}
		}

		close(dataChan)
		close(errChan)
	}()

	return dataChan, errChan
}

func (sl *streamLoader) load(ctx context.Context) (data []any, next bool, err error) {
	rows, err := sl.db.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM %s OFFSET %d LIMIT %d", sl.table, sl.currOffset, sl.queryLimit))
	if err != nil {
		err := fmt.Errorf("error loading data: %w", err)
		return nil, false, err
	}

	for rows.Next() {
		item := make(map[string]any)
		if err := rows.MapScan(item); err != nil {
			err := fmt.Errorf("error scanning row: %w", err)
			return nil, false, err
		}
		data = append(data, item)
		sl.currOffset++

		if sl.sampleLimit > 0 && sl.currOffset >= sl.sampleLimit {
			return data, false, nil
		}
	}

	if len(data) == 0 {
		return data, false, nil
	}

	return data, true, nil
}
