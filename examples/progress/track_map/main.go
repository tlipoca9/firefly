package main

import (
	"fmt"
	"time"

	"github.com/tlipoca9/firefly/progress"
)

func main() {
	data := make(map[string]int)
	for i := 0; i < 42; i++ {
		data["key "+fmt.Sprint(i)] = i
	}

	var newData []string
	progress.TrackMap(data, func(k string, v int) bool {
		newData = append(newData, fmt.Sprintf("%s: %d", k, v*v))
		time.Sleep(200 * time.Millisecond)
		return true
	})

	for _, v := range newData {
		fmt.Println(v)
	}

	time.Sleep(time.Second)
}
