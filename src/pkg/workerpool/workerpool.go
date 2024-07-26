package workerpool

import (
	"fmt"
	"log/slog"
)

type Task struct {
	ID         string
	Payload    any
	Executor   func() (any, error)
	OmitResult bool
}

type TaskResult struct {
	ID       string
	WorkerID string
	Result   any
	Err      error
}

type Worker struct {
	ID        string
	results   chan TaskResult
	taskQueue chan Task
}

func (w *Worker) Start() {
	go func() {
		for task := range w.taskQueue {
			slog.Debug("worker executor", slog.Any("workerID", w.ID), slog.Any("taskID", task.ID))
			result, err := task.Executor()
			if err != nil {
				slog.Error("worker executor", slog.Any("workerID", w.ID), slog.Any("taskID", task.ID), slog.Any("err", err))
			}

			r := TaskResult{WorkerID: w.ID, Result: result, Err: err}

			if !task.OmitResult {
				w.results <- r
			}
		}
	}()
}

type WorkerPool struct {
	taskQueue    chan Task
	results      chan TaskResult
	numOfWorkers int
}

func New(numOfWorkers int) *WorkerPool {
	return &WorkerPool{
		numOfWorkers: numOfWorkers,
		taskQueue:    make(chan Task),
		results:      make(chan TaskResult),
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numOfWorkers; i++ {
		w := &Worker{
			ID:        fmt.Sprintf("%d", i+1),
			results:   wp.results,
			taskQueue: wp.taskQueue,
		}
		w.Start()
	}
}

func (wp *WorkerPool) Submit(task Task) {
	wp.taskQueue <- task
}

func (wp *WorkerPool) GetResult() TaskResult {
	return <-wp.results
}
