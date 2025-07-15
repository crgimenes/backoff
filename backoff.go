package backoff

import (
	"math/rand"
	"sync"
	"time"
)

// Backoff encapsulates the state for exponential backoff.
type Backoff struct {
	mu          sync.Mutex    // garante segurança em concorrência
	initial     time.Duration // valor base
	factor      float64       // fator ≥ 1.0
	max         time.Duration // limite superior
	withJitter  bool          // habilita jitter
	current     time.Duration // último intervalo retornado
	initialized bool          // indica primeira chamada
}

// New cria um Backoff com jitter opcional (default true).
func New(initial time.Duration, factor float64, max time.Duration, opts ...Option) *Backoff {
	b := &Backoff{
		initial:    initial,
		factor:     factor,
		max:        max,
		withJitter: true,
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

// Option permite customizar Backoff.
type Option func(*Backoff)

// WithJitter desabilita ou habilita o jitter.
func WithJitter(enabled bool) Option {
	return func(b *Backoff) {
		b.withJitter = enabled
	}
}

// Next retorna o próximo intervalo, aplicando fator e jitter (se habilitado).
func (b *Backoff) Next() time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()

	// primeira chamada
	if !b.initialized {
		b.current = b.initial
		b.initialized = true
	} else {
		// calcula expoencial
		next := time.Duration(float64(b.current) * b.factor)
		if next > b.max {
			next = b.max
		}
		b.current = next
	}

	// aplica jitter completo: [0, current)
	if b.withJitter {
		return time.Duration(rand.Int63n(int64(b.current + 1)))
	}
	return b.current
}

// Reset reinicia o estado para a primeira chamada.
func (b *Backoff) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.initialized = false
}

// Example of usage:
//
//   b := backoff.New(500*time.Millisecond, 2.0, 10*time.Second)
//   for i := 0; i < 5; i++ {
//       wait := b.Next()
//       fmt.Printf("retry #%d: waiting %v\n", i+1, wait)
//       time.Sleep(wait)
//   }
//
// Para desabilitar jitter:
//
//   b2 := backoff.New(1*time.Second, 2.0, 30*time.Second, backoff.WithJitter(false))
