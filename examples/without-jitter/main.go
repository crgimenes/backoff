package main

import (
	"fmt"
	"time"

	"github.com/crgimenes/backoff"
)

func main() {
	// Disable jitter for predictable delays
	b := backoff.New(1*time.Second, 2.0, 30*time.Second, backoff.WithJitter(false))

	for i := 0; i < 4; i++ {
		wait := b.Next()
		fmt.Printf("Retry #%d: waiting %v\n", i+1, wait)
		time.Sleep(wait)
	}
}
