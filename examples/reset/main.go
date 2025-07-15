package main

import (
	"fmt"
	"time"

	"github.com/crgimenes/backoff"
)

func main() {
	b := backoff.New(100*time.Millisecond, 2.0, 1*time.Second, backoff.WithJitter(false))

	fmt.Println("First operation:")
	// Use for first operation
	fmt.Printf("Call 1: %v\n", b.Next()) // 100ms
	fmt.Printf("Call 2: %v\n", b.Next()) // 200ms
	fmt.Printf("Call 3: %v\n", b.Next()) // 400ms

	fmt.Println("\nAfter reset:")
	// Reset for new operation
	b.Reset()
	fmt.Printf("Call 1: %v (back to initial)\n", b.Next()) // 100ms (back to initial)
	fmt.Printf("Call 2: %v\n", b.Next())                   // 200ms
	fmt.Printf("Call 3: %v\n", b.Next())                   // 400ms

	fmt.Println("\nDemonstrating multiple operations with reset:")
	for operation := 1; operation <= 3; operation++ {
		fmt.Printf("\nOperation #%d:\n", operation)
		b.Reset()
		for attempt := 1; attempt <= 3; attempt++ {
			wait := b.Next()
			fmt.Printf("  Attempt %d: %v\n", attempt, wait)
		}
	}
}
