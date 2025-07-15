package backoff

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		initial    time.Duration
		factor     float64
		max        time.Duration
		opts       []Option
		wantJitter bool
	}{
		{
			name:       "default with jitter enabled",
			initial:    100 * time.Millisecond,
			factor:     2.0,
			max:        5 * time.Second,
			opts:       nil,
			wantJitter: true,
		},
		{
			name:       "with jitter disabled",
			initial:    200 * time.Millisecond,
			factor:     1.5,
			max:        10 * time.Second,
			opts:       []Option{WithJitter(false)},
			wantJitter: false,
		},
		{
			name:       "with jitter explicitly enabled",
			initial:    50 * time.Millisecond,
			factor:     3.0,
			max:        1 * time.Second,
			opts:       []Option{WithJitter(true)},
			wantJitter: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(tt.initial, tt.factor, tt.max, tt.opts...)
			if b.initial != tt.initial {
				t.Errorf("New() initial = %v, want %v", b.initial, tt.initial)
			}
			if b.factor != tt.factor {
				t.Errorf("New() factor = %v, want %v", b.factor, tt.factor)
			}
			if b.max != tt.max {
				t.Errorf("New() max = %v, want %v", b.max, tt.max)
			}
			if b.withJitter != tt.wantJitter {
				t.Errorf("New() withJitter = %v, want %v", b.withJitter, tt.wantJitter)
			}
			if b.initialized {
				t.Errorf("New() initialized should be false initially")
			}
		})
	}
}

func TestWithJitter(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
		want    bool
	}{
		{"enable jitter", true, true},
		{"disable jitter", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backoff{}
			opt := WithJitter(tt.enabled)
			opt(b)
			if b.withJitter != tt.want {
				t.Errorf("WithJitter() = %v, want %v", b.withJitter, tt.want)
			}
		})
	}
}

func TestBackoff_Next(t *testing.T) {
	tests := []struct {
		name    string
		initial time.Duration
		factor  float64
		max     time.Duration
		jitter  bool
		calls   int
		checkFn func(t *testing.T, b *Backoff, results []time.Duration)
	}{
		{
			name:    "without jitter - exponential growth",
			initial: 100 * time.Millisecond,
			factor:  2.0,
			max:     1 * time.Second,
			jitter:  false,
			calls:   4,
			checkFn: func(t *testing.T, b *Backoff, results []time.Duration) {
				expected := []time.Duration{
					100 * time.Millisecond,
					200 * time.Millisecond,
					400 * time.Millisecond,
					800 * time.Millisecond,
				}
				for i, want := range expected {
					if results[i] != want {
						t.Errorf("Next() call %d = %v, want %v", i+1, results[i], want)
					}
				}
			},
		},
		{
			name:    "without jitter - hits max limit",
			initial: 500 * time.Millisecond,
			factor:  2.0,
			max:     800 * time.Millisecond,
			jitter:  false,
			calls:   3,
			checkFn: func(t *testing.T, b *Backoff, results []time.Duration) {
				expected := []time.Duration{
					500 * time.Millisecond,
					800 * time.Millisecond, // capped at max
					800 * time.Millisecond, // stays at max
				}
				for i, want := range expected {
					if results[i] != want {
						t.Errorf("Next() call %d = %v, want %v", i+1, results[i], want)
					}
				}
			},
		},
		{
			name:    "with jitter - values within range",
			initial: 100 * time.Millisecond,
			factor:  2.0,
			max:     1 * time.Second,
			jitter:  true,
			calls:   3,
			checkFn: func(t *testing.T, b *Backoff, results []time.Duration) {
				// First call should be [0, 100ms]
				if results[0] < 0 || results[0] > 100*time.Millisecond {
					t.Errorf("Next() call 1 = %v, want range [0, %v]", results[0], 100*time.Millisecond)
				}
				// Second call should be [0, 200ms]
				if results[1] < 0 || results[1] > 200*time.Millisecond {
					t.Errorf("Next() call 2 = %v, want range [0, %v]", results[1], 200*time.Millisecond)
				}
				// Third call should be [0, 400ms]
				if results[2] < 0 || results[2] > 400*time.Millisecond {
					t.Errorf("Next() call 3 = %v, want range [0, %v]", results[2], 400*time.Millisecond)
				}
			},
		},
		{
			name:    "factor of 1.0 - no growth",
			initial: 200 * time.Millisecond,
			factor:  1.0,
			max:     1 * time.Second,
			jitter:  false,
			calls:   3,
			checkFn: func(t *testing.T, b *Backoff, results []time.Duration) {
				for i, result := range results {
					if result != 200*time.Millisecond {
						t.Errorf("Next() call %d = %v, want %v", i+1, result, 200*time.Millisecond)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set a fixed seed for reproducible jitter tests
			rand.Seed(42)

			b := New(tt.initial, tt.factor, tt.max, WithJitter(tt.jitter))
			var results []time.Duration

			for i := 0; i < tt.calls; i++ {
				results = append(results, b.Next())
			}

			tt.checkFn(t, b, results)
		})
	}
}

func TestBackoff_Reset(t *testing.T) {
	b := New(100*time.Millisecond, 2.0, 1*time.Second, WithJitter(false))

	// Call Next a few times to advance state
	first := b.Next()
	second := b.Next()

	// Verify progression
	if first != 100*time.Millisecond {
		t.Errorf("First Next() = %v, want %v", first, 100*time.Millisecond)
	}
	if second != 200*time.Millisecond {
		t.Errorf("Second Next() = %v, want %v", second, 200*time.Millisecond)
	}

	// Reset and verify it starts over
	b.Reset()
	afterReset := b.Next()
	if afterReset != 100*time.Millisecond {
		t.Errorf("Next() after Reset() = %v, want %v", afterReset, 100*time.Millisecond)
	}
}

func TestBackoff_Concurrency(t *testing.T) {
	b := New(10*time.Millisecond, 1.5, 100*time.Millisecond, WithJitter(false))

	const numGoroutines = 100
	const callsPerGoroutine = 10

	var wg sync.WaitGroup
	results := make([][]time.Duration, numGoroutines)

	// Launch concurrent goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			for j := 0; j < callsPerGoroutine; j++ {
				duration := b.Next()
				results[index] = append(results[index], duration)
			}
		}(i)
	}

	wg.Wait()

	// Verify all results are valid (non-zero and within reasonable bounds)
	for i, goroutineResults := range results {
		for j, duration := range goroutineResults {
			if duration <= 0 {
				t.Errorf("Goroutine %d, call %d: got non-positive duration %v", i, j, duration)
			}
			if duration > 100*time.Millisecond {
				t.Errorf("Goroutine %d, call %d: got duration %v exceeding max %v", i, j, duration, 100*time.Millisecond)
			}
		}
	}
}

func TestBackoff_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		initial time.Duration
		factor  float64
		max     time.Duration
		jitter  bool
	}{
		{
			name:    "zero initial duration",
			initial: 0,
			factor:  2.0,
			max:     1 * time.Second,
			jitter:  false,
		},
		{
			name:    "very small factor",
			initial: 100 * time.Millisecond,
			factor:  1.001,
			max:     1 * time.Second,
			jitter:  false,
		},
		{
			name:    "max smaller than initial",
			initial: 1 * time.Second,
			factor:  2.0,
			max:     500 * time.Millisecond,
			jitter:  false,
		},
		{
			name:    "very large factor",
			initial: 1 * time.Nanosecond,
			factor:  1000.0,
			max:     1 * time.Second,
			jitter:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := New(tt.initial, tt.factor, tt.max, WithJitter(tt.jitter))

			// Should not panic and should return reasonable values
			for i := 0; i < 5; i++ {
				duration := b.Next()
				if duration < 0 {
					t.Errorf("Next() call %d returned negative duration: %v", i+1, duration)
				}
				// For edge case where max < initial, allow first call to return initial
				if tt.name == "max smaller than initial" && i == 0 {
					// First call can return initial value even if it exceeds max
					continue
				}
				if duration > tt.max {
					t.Errorf("Next() call %d returned duration %v exceeding max %v", i+1, duration, tt.max)
				}
			}
		})
	}
}

func BenchmarkBackoff_Next(b *testing.B) {
	backoff := New(100*time.Millisecond, 2.0, 10*time.Second, WithJitter(false))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = backoff.Next()
	}
}

func BenchmarkBackoff_NextWithJitter(b *testing.B) {
	backoff := New(100*time.Millisecond, 2.0, 10*time.Second, WithJitter(true))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = backoff.Next()
	}
}

func BenchmarkBackoff_Concurrent(b *testing.B) {
	backoff := New(100*time.Millisecond, 2.0, 10*time.Second, WithJitter(false))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = backoff.Next()
		}
	})
}
