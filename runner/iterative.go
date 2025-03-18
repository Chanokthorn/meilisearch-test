package runner

import (
	"context"
	"ms-tester/meilisearch"
	"ms-tester/storage"
)

type Iterative struct {
	ms meilisearch.MeiliSearch
	IterativeConfig
	
}

type IterativeConfig struct {
	storage 	storage.Storage 
	// dataStream   chan any
	indexUid     string
	// workerAmount int
}

func NewIterative(ms meilisearch.MeiliSearch) *Iterative {
	return &Iterative{ms: ms}
}

// func (i *Iterative) SetData(data []any) *Iterative {
// 	i.data = data
// 	return i
// }

// func (i *Iterative) SetDataStream(dataStream chan any) *Iterative {
// 	i.dataStream = dataStream
// 	return i
// }

func (i *Iterative) SetStorage(storage storage.Storage) *Iterative {
	i.storage = storage
	return i
}

func (i *Iterative) SetIndexUid(indexUid string) *Iterative {
	i.indexUid = indexUid
	return i
}

// func (i *Iterative) SetWorkerAmount(workerAmount int) *Iterative {
// 	i.workerAmount = workerAmount
// 	return i
// }

func (i *Iterative) Run(ctx context.Context) (int, error) {
	var (
		lastTaskUid int
		err         error
	)

	for _, d := range i.storage.

	for _, d := range i.data {
		lastTaskUid, err = i.ms.AddOrUpdateDocument(ctx, i.indexUid, d)
		if err != nil {
			return 0, err
		}
	}

	return lastTaskUid, nil
}

// func (i *Iterative) RunStream(ctx context.Context) (int, error) {
// 	var (
// 		taskUidChan = make(chan int)
// 		errChan      = make(chan error)
// 	)

// 	for x := 0; x < i.workerAmount; x++ {
// 		go func() {
// 			for d := range i.dataStream {
// 				taskUid, err := i.ms.AddOrUpdateDocument(ctx, i.indexUid, d)
// 				if err != nil {
// 					errChan <- err
// 					return
// 				}

// 				taskUidChan <- taskUid
// 			}
// 		}()
// 	}

// 	var lastTaskUid int
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return 0, ctx.Err()
// 		case err := <-errChan:
// 			return 0, err
// 		case lastTaskUid = <-taskUidChan:
// 		}
// 	}
// }
