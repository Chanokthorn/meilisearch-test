package runner

import (
	"context"
	"fmt"
	"log"
	"ms-tester/meilisearch"
	"ms-tester/storage"
	"sync"
)

type Iterative struct {
	ms           meilisearch.MeiliSearch
	workerAmount int
	IterativeConfig
	mu sync.Mutex
}

type IterativeConfig struct {
	loader   storage.StreamLoader
	indexUid string
}

func NewIterative(ms meilisearch.MeiliSearch) *Iterative {
	return &Iterative{
		workerAmount: 10,
		ms: ms,
	}
}

func (i *Iterative) SetIndexUid(indexUid string) *Iterative {
	i.indexUid = indexUid
	return i
}

func (i *Iterative) SetWorkerAmount(workerAmount int) *Iterative {
	i.workerAmount = workerAmount
	return i
}

func (i *Iterative) SetLoader(loader storage.StreamLoader) *Iterative {
	i.loader = loader
	return i
}

func (i *Iterative) Run(ctx context.Context) (int, error) {
	if i.loader == nil {
		return 0, fmt.Errorf("loader is not set")
	}

	var (
		lastTaskUidChan = make(chan int)
		lastTaskUid     int
		uploadErrChan   = make(chan error, 10)
	)

	// run stream loader
	dataChan, _ := i.loader.Start()

	var workerWG sync.WaitGroup

	// for worker in worker amount, run function that consumes from stream loader
	for x := 0; x < i.workerAmount; x++ {
		workerWG.Add(1)
		go func() {
			defer workerWG.Done()
			for d := range dataChan {
				taskUid, err := i.ms.AddOrUpdateDocument(ctx, i.indexUid, d)
				if err != nil {
					uploadErrChan <- err
					return
				}

				lastTaskUidChan <- taskUid
			}
		}()
	}

	// collect information from workers
	var collectorWG sync.WaitGroup

	collectorWG.Add(1)
	go func() {
		defer collectorWG.Done()
		for err := range uploadErrChan {
			err = fmt.Errorf("error uploading document: %w", err)
			log.Println(err.Error())
			return
		}
	}()

	collectorWG.Add(1)
	go func() {
		defer collectorWG.Done()
		for taskUid := range lastTaskUidChan {
			i.mu.Lock()
			lastTaskUid = max(lastTaskUid, taskUid)
			i.mu.Unlock()
		}
	}()

	workerWG.Wait()
	close(lastTaskUidChan)
	close(uploadErrChan)

	collectorWG.Wait()

	return lastTaskUid, nil
}
