# Backoff

A thread-safe, configurable exponential backoff implementation for Go applications. This package provides a simple and efficient way to implement retry logic with exponential delays, commonly used in distributed systems, API clients, and error recovery mechanisms.

## Features

- **Thread-safe**: Safe for concurrent use across multiple goroutines
- **Configurable**: Customizable initial delay, growth factor, and maximum delay
- **Optional Jitter**: Helps prevent thundering herd problems in distributed systems
- **Zero Dependencies**: Uses only Go's standard library
- **Lightweight**: Minimal memory footprint and CPU overhead

## Installation

```bash
go get github.com/crgimenes/backoff
```

## Usage

### Basic Usage

```go
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
```

**Expected Output:**

```
Retry #1: waiting 245ms    // Random jitter applied: [0, 500ms]
Retry #2: waiting 891ms    // Random jitter applied: [0, 1000ms]
Retry #3: waiting 1.2s     // Random jitter applied: [0, 2000ms]
Retry #4: waiting 3.1s     // Random jitter applied: [0, 4000ms]
Retry #5: waiting 7.8s     // Random jitter applied: [0, 8000ms]
```

### Without Jitter

```go
// Disable jitter for predictable delays
b := backoff.New(1*time.Second, 2.0, 30*time.Second, backoff.WithJitter(false))

for i := 0; i < 4; i++ {
    wait := b.Next()
    fmt.Printf("Retry #%d: waiting %v\n", i+1, wait)
}
```

**Expected Output:**

```
Retry #1: waiting 1s
Retry #2: waiting 2s
Retry #3: waiting 4s
Retry #4: waiting 8s
```

### Retry Pattern with HTTP Client

```go
package main

import (
    "fmt"
    "net/http"
    "time"

    "github.com/crgimenes/backoff"
)

func fetchWithRetry(url string, maxRetries int) (*http.Response, error) {
    b := backoff.New(100*time.Millisecond, 2.0, 5*time.Second)

    var resp *http.Response
    var err error

    for i := 0; i < maxRetries; i++ {
        resp, err = http.Get(url)
        if err == nil && resp.StatusCode == http.StatusOK {
            return resp, nil
        }

        if i < maxRetries-1 {
            wait := b.Next()
            fmt.Printf("Request failed, retrying in %v...\n", wait)
            time.Sleep(wait)
        }
    }

    return resp, err
}
```

### Reset Functionality

```go
b := backoff.New(100*time.Millisecond, 2.0, 1*time.Second)

// Use for first operation
b.Next() // 100ms
b.Next() // 200ms
b.Next() // 400ms

// Reset for new operation
b.Reset()
b.Next() // 100ms (back to initial)
```

## API Reference

### Types

#### `Backoff`

Thread-safe exponential backoff calculator.

#### `Option`

Configuration function type for customizing Backoff behavior.

### Functions

#### `New(initial time.Duration, factor float64, max time.Duration, opts ...Option) *Backoff`

Creates a new Backoff instance.

**Parameters:**

- `initial`: Starting delay duration
- `factor`: Multiplier for each subsequent delay (should be ≥ 1.0)
- `max`: Maximum delay duration (upper bound)
- `opts`: Optional configuration functions

**Returns:** Configured Backoff instance

#### `WithJitter(enabled bool) Option`

Option to enable or disable jitter.

**Parameters:**

- `enabled`: true to enable jitter (default), false to disable

**Returns:** Option function

### Methods

#### `(b *Backoff) Next() time.Duration`

Calculates and returns the next delay duration.

**Returns:** Duration to wait before next retry

#### `(b *Backoff) Reset()`

Resets the backoff state to initial values.

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

## Performance

The backoff implementation is highly optimized for performance:

- **Memory**: Minimal allocations, reuses internal state
- **CPU**: Fast O(1) calculations with simple arithmetic
- **Concurrency**: Efficient mutex-based synchronization

Benchmark results on typical hardware:

```
BenchmarkBackoff_Next-8                50000000    25.4 ns/op    0 B/op    0 allocs/op
BenchmarkBackoff_NextWithJitter-8      30000000    45.2 ns/op    0 B/op    0 allocs/op
BenchmarkBackoff_Concurrent-8          20000000    67.8 ns/op    0 B/op    0 allocs/op
```

## Best Practices

1. **Choose appropriate parameters**: Start with reasonable initial delays (100ms-1s) and factors (1.5-2.0)
2. **Set reasonable maximums**: Prevent excessive delays with appropriate max values (10s-60s)
3. **Use jitter in distributed systems**: Helps prevent thundering herd effects
4. **Reset between operations**: Call `Reset()` when starting new retry sequences
5. **Consider context cancellation**: Combine with `context.Context` for proper timeout handling

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Run `go fmt` to format code
7. Commit your changes (`git commit -am 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Development Guidelines

- Follow Go conventions and idioms
- Write table-driven tests for new features
- Maintain backward compatibility
- Update documentation for API changes
- Ensure cross-platform compatibility (Linux, macOS, Windows)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Related Projects

- [cenkalti/backoff](https://github.com/cenkalti/backoff) - More feature-rich backoff library
- [jpillora/backoff](https://github.com/jpillora/backoff) - Simple backoff implementation
- [lestrrat-go/backoff](https://github.com/lestrrat-go/backoff) - Highly configurable backoff library

## Changelog

### v1.0.0

- Initial release
- Thread-safe exponential backoff
- Optional jitter support
- Zero-dependency implementation
