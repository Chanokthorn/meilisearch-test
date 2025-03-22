package runner

import (
	"context"
	"fmt"
	"ms-tester/meilisearch"
)

type iterativeWorker struct {
	ms        meilisearch.MeiliSearch
	indexUid  string
}

func NewIterativeWorker(ms meilisearch.MeiliSearch) *iterativeWorker {
	return &iterativeWorker{
		ms:        ms,
	}
}

func (iw *iterativeWorker) SetIndexUid(uid string) *iterativeWorker {
	iw.indexUid = uid
	return iw
}

// process for iterative worker
func (iw *iterativeWorker) Process(ctx context.Context, dataChan <-chan any, taskIDChan chan<- int, errChan chan<- error) {
	for d := range dataChan {
		taskID, err := iw.ms.AddOrUpdateDocument(ctx, iw.indexUid, d)
		if err != nil {
			err := fmt.Errorf("error uploading document: %w", err)
			errChan <- err
		}

		taskIDChan <- taskID
	}
	
}