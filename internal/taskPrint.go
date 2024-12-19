package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PrintTask struct {
	id      string
	message string
	owner   string
}

func NewPrintTask(message string) *PrintTask {
	return &PrintTask{
		id:      uuid.New().String(),
		message: message,
	}
}

func (pt *PrintTask) Execute(ctx context.Context) error {
	fmt.Printf("[%v] %v ", pt.owner, pt.message)
	return nil
}

func (pt *PrintTask) GetID() string {
	return pt.id
}

func (pt *PrintTask) SetOwner(own string) {
	pt.owner = own
}

type TaskOption func(p *PrintTask) error

func CreatePrintTasks(interval time.Duration, opts ...TaskOption) error {

	ctx, cancel := context.WithTimeout(context.Background(), interval)
	defer cancel()

	tq := NewTaskQueue()

	wr := NewWorker(tq, "1")
	wr.Run(ctx)

	wr2 := NewWorker(tq, "2")
	wr2.Run(ctx)

	err := tq.Start(ctx)
	if err != nil {
		panic(err)
	}

	tiker := time.NewTicker(time.Second)
	for t := range tiker.C {
		pt := NewPrintTask(fmt.Sprintf("%v", t.Format(time.RFC3339)))

		for _, opt := range opts {
			if err := opt(pt); err != nil {
				fmt.Println(err)
			}
		}
		if err := tq.Add(pt); err != nil {
			return nil
		}
	}

	go func() {
		<-ctx.Done()
		tiker.Stop()
	}()

	return nil
}
