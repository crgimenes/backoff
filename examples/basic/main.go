package main

import (
	"fmt"
	"time"

	"github.com/crgimenes/backoff"
)

func main() {
	// Create a backoff with 500ms initial delay, 2x growth factor, 10s max delay
	b := backoff.New(500*time.Millisecond, 2.0, 10*time.Second)

	for i := range 5 {
		wait := b.Next()
		fmt.Printf("Retry #%d: waiting %v\n", i+1, wait)
		time.Sleep(wait)
	}
}
