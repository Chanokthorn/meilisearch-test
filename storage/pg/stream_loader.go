package pg

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type streamLoader struct {
	db *sqlx.DB
}

func NewStreamLoader(url string) (*streamLoader, error) {
	pgConn, err := sqlx.Connect("postgres", url)
	if err != nil {
		err := fmt.Errorf("unable to connect to postgres: %w", err)
		return nil, err
	}

	return &streamLoader{db: pgConn}, nil
}

func (sl *streamLoader) Start() (chan interface{}, chan error) {
}
