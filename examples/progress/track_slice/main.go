package main

import (
	"fmt"
	"time"

	"github.com/tlipoca9/firefly/progress"
)

func main() {
	var data []int
	for i := 0; i < 42; i++ {
		data = append(data, i)
	}

	var newData []int
	progress.TrackSlice(data, func(i int, v int) bool {
		newData = append(newData, v*v)
		time.Sleep(200 * time.Millisecond)
		return true
	})

	for _, v := range newData {
		fmt.Print(v)
		fmt.Print(" ")
	}
	fmt.Println()

	time.Sleep(time.Second)
}
