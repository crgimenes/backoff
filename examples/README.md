# Examples

This directory contains practical examples demonstrating different use cases of the backoff library.

## Running Examples

Each example is in its own subdirectory. To run an example:

```bash
cd examples/<example-name>
go run main.go
```

## Available Examples

### 1. Basic Usage (`basic/`)

Demonstrates the basic exponential backoff with jitter enabled (default behavior).

```bash
cd examples/basic
go run main.go
```

**What it shows:**

- Default jitter behavior
- Exponential growth with randomization
- Basic retry loop pattern

### 2. Without Jitter (`without-jitter/`)

Shows predictable, deterministic backoff delays without randomization.

```bash
cd examples/without-jitter
go run main.go
```

**What it shows:**

- Disabling jitter for predictable delays
- Pure exponential growth: 1s → 2s → 4s → 8s
- Use case: When you need deterministic timing

### 3. HTTP Retry Pattern (`http-retry/`)

Real-world example of retrying HTTP requests with exponential backoff.

```bash
cd examples/http-retry
go run main.go
```

**What it shows:**

- Practical retry pattern for HTTP clients
- Error handling and response validation
- Exponential backoff between failed requests
- Proper resource cleanup

### 4. Reset Functionality (`reset/`)

Demonstrates how to reset backoff state for new operations.

```bash
cd examples/reset
go run main.go
```

**What it shows:**

- Using `Reset()` to restart backoff sequences
- Multiple independent operations
- State management between different retry scenarios

## Example Output

### Basic Example

```
Retry #1: waiting 245ms
Retry #2: waiting 891ms
Retry #3: waiting 1.2s
Retry #4: waiting 3.1s
Retry #5: waiting 7.8s
```

### Without Jitter Example

```
Retry #1: waiting 1s
Retry #2: waiting 2s
Retry #3: waiting 4s
Retry #4: waiting 8s
```

### HTTP Retry Example

```
Attempting to fetch https://httpbin.org/status/500 with 5 retries...
Request failed, retrying in 67ms...
Request failed, retrying in 134ms...
Request failed, retrying in 268ms...
Request failed, retrying in 536ms...
Final error: <error details>
```

### Reset Example

```
First operation:
Call 1: 100ms
Call 2: 200ms
Call 3: 400ms

After reset:
Call 1: 100ms (back to initial)
Call 2: 200ms
Call 3: 400ms

Demonstrating multiple operations with reset:

Operation #1:
  Attempt 1: 100ms
  Attempt 2: 200ms
  Attempt 3: 400ms

Operation #2:
  Attempt 1: 100ms
  Attempt 2: 200ms
  Attempt 3: 400ms

Operation #3:
  Attempt 1: 100ms
  Attempt 2: 200ms
  Attempt 3: 400ms
```

## Common Patterns

These examples demonstrate common patterns you can use in your applications:

1. **API Clients**: Use the HTTP retry pattern for robust API communication
2. **Database Operations**: Apply reset functionality for independent transaction retries
3. **Message Processing**: Use basic backoff for handling temporary processing failures
4. **Testing**: Use without-jitter for predictable test scenarios

## Customization

You can modify these examples to experiment with different parameters:

- **Initial delay**: Start with different base delays (100ms, 1s, etc.)
- **Growth factor**: Try different multipliers (1.5, 2.0, 3.0)
- **Maximum delay**: Set different upper bounds (5s, 30s, 2m)
- **Jitter**: Enable/disable based on your needs
