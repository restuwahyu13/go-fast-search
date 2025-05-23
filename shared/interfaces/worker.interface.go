package inf

import "sync"

type (
	ISearchWorker interface {
		SearchRun(wg *sync.WaitGroup)
	}

	IDeadLetterQueueWorker interface {
		DeadLetterQueueRun(wg *sync.WaitGroup)
	}
)
