package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/tlipoca9/firefly/progress"
)

func main() {
	client := http.Client{
		Transport: progress.TrackHTTP("Speed testing", http.DefaultTransport),
	}

	req, err := http.NewRequest(
		http.MethodGet,
		"https://speed.cloudflare.com/__down?during=download&bytes=104857600",
		nil,
	)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Write the response body to file
	f, err := os.OpenFile("speed_test.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second)
}
