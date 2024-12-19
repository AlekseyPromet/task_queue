package main

import (
	"AlekseyPromet/algo/taskqueue/internal"
	"time"
)

func main() {

	optTask := func(p *internal.PrintTask) error {
		p.SetOwner("root")
		return nil
	}

	if err := internal.CreatePrintTasks(time.Second*10, []internal.TaskOption{optTask}...); err != nil {
		panic(err)
	}
}
