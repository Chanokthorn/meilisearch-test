package runner

import (
	"context"
	"fmt"
	"log"
	"ms-tester/storage"
	"sync"
)

type Runner interface {
	SetWorkerAmount(workerAmount int) Runner
	SetWorker(worker Worker) Runner
	Run(ctx context.Context, loader storage.StreamLoader) (finalTaskUID int, err error)
}

type Worker interface {
	Process(ctx context.Context, dataChan <-chan any, taskIDChan chan<- int, errChan chan<- error)
}

type runner struct {
	worker       Worker
	workerAmount int
	mu           sync.Mutex
}

func NewRunner() Runner {
	return &runner{
		workerAmount: 10,
	}
}

func (i *runner) SetWorkerAmount(workerAmount int) Runner {
	i.workerAmount = workerAmount
	return i
}

func (br *runner) SetWorker(w Worker) Runner {
	br.worker = w
	return br
}

func (rn *runner) Run(ctx context.Context, loader storage.StreamLoader) (int, error) {
	if rn.worker == nil {
		log.Fatal("worker is not set")
	}
	var (
		lastTaskUidChan = make(chan int)
		lastTaskUid     int
		uploadErrChan   = make(chan error, 10)
	)

	// run stream loader
	dataChan, _ := loader.Start()

	var workerWG sync.WaitGroup

	// for worker in worker amount, run function that consumes from stream loader
	for x := 0; x < rn.workerAmount; x++ {
		workerWG.Add(1)
		go func() {
			defer workerWG.Done()
			rn.worker.Process(ctx, dataChan, lastTaskUidChan, uploadErrChan)
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
			rn.mu.Lock()
			lastTaskUid = max(lastTaskUid, taskUid)
			rn.mu.Unlock()
		}
	}()

	workerWG.Wait()
	close(lastTaskUidChan)
	close(uploadErrChan)

	collectorWG.Wait()

	return lastTaskUid, nil
}
