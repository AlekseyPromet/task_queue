package internal

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Task interface {
	GetID() string
	Execute(ctx context.Context) error
	SetOwner(owner string)
}

type TaskQueue struct {
	isActive bool
	tasks    []Task
	enqueue  chan Task
	dequeue  chan Task
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		isActive: true,
		tasks:    make([]Task, 0),
		enqueue:  make(chan Task),
		dequeue:  make(chan Task),
	}
}

type Worker struct {
	id        string
	taskQueue *TaskQueue
}

func NewWorker(tq *TaskQueue, id string) *Worker {
	wr := &Worker{
		id:        id,
		taskQueue: tq,
	}

	if id == "" {
		wr.id = uuid.NewString()
	}

	return wr
}

func (w *Worker) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case task := <-w.taskQueue.dequeue:
				if err := task.Execute(ctx); err != nil {
					// Handle error
					fmt.Printf("[ERROR][worker %s] Task %s is executing failed\n", w.id, task.GetID())
					continue
				}
				fmt.Printf("[worker %s] Task %s is executed\n", w.id, task.GetID())
			}
		}
	}()
}

func (tq *TaskQueue) Start(ctx context.Context) error {

	go func() {
		for {
			select {
			case <-ctx.Done():
				tq.isActive = false
				close(tq.enqueue)
				return
			case task := <-tq.enqueue:
				tq.tasks = append(tq.tasks, task)
			default:
				if len(tq.tasks) > 0 {
					tq.dequeue <- tq.tasks[0]
					tq.tasks = tq.tasks[1:]
				}
			}
		}
	}()

	return nil
}

func (tq *TaskQueue) Add(task Task) error {
	if tq.isActive {
		tq.enqueue <- task
		return nil
	}

	return fmt.Errorf("stop executing tasks")
}
