package main

import (
	"time"

	"github.com/tlipoca9/firefly/progress"
)

func main() {
	p := progress.New(progress.NewTaskModel("Simple", progress.WithTotal(42)))
	p.Start()

	for i := 0; i < 42; i++ {
		p.Incr(1)
		time.Sleep(200 * time.Millisecond)
	}

	p.Finish()

	time.Sleep(time.Second)
}
