package workerpool

import "github.com/alitto/pond"

var _consumer_pool *pond.WorkerPool

func InitPool(maxWorkers, maxCapacity int) {
	_consumer_pool = pond.New(maxWorkers, maxCapacity)
}

func Use() *pond.WorkerPool {
	return _consumer_pool
}

func Close() {
	_consumer_pool.StopAndWait()
}
