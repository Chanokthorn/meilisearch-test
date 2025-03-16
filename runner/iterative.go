package runner

import (
	"context"
	"ms-tester/meilisearch"
)

type Iterative struct {
	ms     meilisearch.MeiliSearch
	IterativeConfig
}

type IterativeConfig struct {
	data     []any
	indexUid string
}

func NewIterative(ms meilisearch.MeiliSearch) *Iterative {
	return &Iterative{ms: ms}
}

func (i *Iterative) SetData(data []any) *Iterative{
	i.data = data
	return i
}

func (i *Iterative) SetIndexUid(indexUid string) *Iterative{
	i.indexUid = indexUid
	return i
}

func (i *Iterative) Run(ctx context.Context) (int, error) {
	var (
		lastTaskUid int
		err         error
	)

	for _, d := range i.data {
		lastTaskUid, err = i.ms.AddOrUpdateDocument(ctx, i.indexUid, d)
		if err != nil {
			return 0, err
		}
	}

	return lastTaskUid, nil
}
