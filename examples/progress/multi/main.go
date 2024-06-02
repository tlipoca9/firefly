package main

import (
	"sync"
	"time"

	"github.com/tlipoca9/firefly/progress"
)

func main() {
	p := progress.New(
		progress.NewTaskModel("Task 1", progress.WithTotal(42)),
		progress.NewTaskModel("Task 2", progress.WithTotal(84)),
		progress.NewTaskModel("Task 3", progress.WithTotal(168)),
	)
	p.Start()

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		for i := 0; i < 42; i++ {
			p.Incr(1, "Task 1")
			time.Sleep(200 * time.Millisecond)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 84; i++ {
			p.Incr(1, "Task 2")
			time.Sleep(200 * time.Millisecond)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 168; i++ {
			p.Incr(1, "Task 3")
			time.Sleep(200 * time.Millisecond)
		}
	}()

	wg.Wait()
	p.Finish()

	time.Sleep(time.Second)
}
