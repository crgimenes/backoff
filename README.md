# Backoff

Exponential backoff implementation for Go with optional jitter and thread safety. Implements retry logic with exponential delays for distributed systems, API clients, and error recovery.

## Features

- **Thread-safe**: Concurrent use across multiple goroutines
- **Configurable**: Initial delay, growth factor, and maximum delay
- **Optional Jitter**: Prevents thundering herd problems
- **Zero Dependencies**: Uses only Go's standard library

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

## Examples

Complete working examples are available in the [`examples/`](examples/) directory:

- **[Basic Usage](examples/basic/)** - Standard exponential backoff with jitter
- **[Without Jitter](examples/without-jitter/)** - Predictable delays without randomization  
- **[HTTP Retry](examples/http-retry/)** - Real-world HTTP client retry pattern
- **[Reset Functionality](examples/reset/)** - Demonstrate state reset between operations

To run any example:

```bash
cd examples/<example-name>
go run main.go
```

Or run all examples:

```bash
cd examples
make all
```

## API Reference

### Types

#### `Backoff`
Exponential backoff calculator.

#### `Option`
Configuration function.

### Functions

#### `New(initial time.Duration, factor float64, max time.Duration, opts ...Option) *Backoff`
Creates a new Backoff instance.

**Parameters:**
- `initial`: Starting delay
- `factor`: Multiplier for each delay (≥ 1.0)
- `max`: Maximum delay
- `opts`: Configuration options

#### `WithJitter(enabled bool) Option`
Enables or disables jitter.

### Methods

#### `(b *Backoff) Next() time.Duration`
Returns the next delay duration.

#### `(b *Backoff) Reset()`
Resets the backoff state.

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

O(1) calculations with zero allocations per call. Uses mutex for thread safety.

Benchmark results:
```
BenchmarkBackoff_Next-8                50000000    25.4 ns/op    0 B/op    0 allocs/op
BenchmarkBackoff_NextWithJitter-8      30000000    45.2 ns/op    0 B/op    0 allocs/op
BenchmarkBackoff_Concurrent-8          20000000    67.8 ns/op    0 B/op    0 allocs/op
```

## Best Practices

1. **Parameters**: Use initial delays of 100ms-1s and factors of 1.5-2.0
2. **Maximum delays**: Set max values between 10s-60s to prevent excessive waits
3. **Jitter**: Enable in distributed systems to prevent thundering herd
4. **Reset**: Call `Reset()` when starting new retry sequences
5. **Timeouts**: Combine with `context.Context` for timeout handling

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/name`)
3. Add tests for new functionality
4. Ensure all tests pass (`go test ./...`)
5. Run `go fmt` to format code
6. Commit and push changes
7. Open a Pull Request

### Guidelines

- Follow Go conventions
- Write table-driven tests
- Maintain backward compatibility
- Update documentation for API changes

## License

MIT License - see [LICENSE](LICENSE) file.

## Related Projects

- [cenkalti/backoff](https://github.com/cenkalti/backoff)
- [jpillora/backoff](https://github.com/jpillora/backoff)
- [lestrrat-go/backoff](https://github.com/lestrrat-go/backoff)

## Changelog

### v1.0.0

- Initial release
- Thread-safe exponential backoff
- Optional jitter support
- Zero-dependency implementation
