package runner

import (
	"context"
	"fmt"
	"ms-tester/meilisearch"
)



type batchWorker struct {
	ms        meilisearch.MeiliSearch
	indexUid  string
	batchSize int
}

func NewBatchWorker(ms meilisearch.MeiliSearch) *batchWorker {
	return &batchWorker{
		ms:        ms,
		batchSize: 100,
	}
}

func (bw *batchWorker) SetIndexUid(uid string) *batchWorker {
	bw.indexUid = uid
	return bw
}

func (bw *batchWorker) SetBatchSize(batchSize int) *batchWorker {
	bw.batchSize = batchSize
	return bw
}

func (bw *batchWorker) Process(ctx context.Context, dataChan <-chan any, taskIDChan chan<- int, errChan chan<- error) {
	var batch []any

	for {
		data, ok := <-dataChan
		if ok {
			batch = append(batch, data)
		}
		// collect data in batch unitl batch size hits limit or dataChan closes
		if len(batch) >= bw.batchSize || !ok {
			taskID, err := bw.ms.AddOrUpdateDocument(ctx, bw.indexUid, batch)
			if err != nil {
				err := fmt.Errorf("unable to upload to ms: %w", err)
				errChan <- err
			}
			// reset batch
			batch = []any{}
			taskIDChan <- taskID
		}

		if !ok {
			break
		}
	}

}