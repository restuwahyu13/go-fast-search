package inf

type (
	ISearchWorker interface {
		SearchRun()
	}

	IDeadLetterQueueWorker interface {
		DeadLetterQueueRun()
	}
)
