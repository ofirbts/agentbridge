package failure

import (
	"errors"
	"time"
)

type RetryPolicy struct {
	MaxRetries int
	BackoffMs  int
	retries    int
}

func NewRetryPolicy(max int, backoff int) RetryPolicy {
	return RetryPolicy{MaxRetries: max, BackoffMs: backoff}
}

func (p *RetryPolicy) RetriesDone() int64 {
	return int64(p.retries)
}

func (p *RetryPolicy) Reset() {
	p.retries = 0
}

func RetryWithPolicy(policy *RetryPolicy, fn func() error) error {
	var lastErr error
	backoff := policy.BackoffMs
	for i := 0; i <= policy.MaxRetries; i++ {
		if i > 0 {
			time.Sleep(time.Duration(backoff) * time.Millisecond)
			policy.retries++
			backoff *= 2
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
	}
	return lastErr
}

var ErrSimulated = errors.New("simulated failure")

func SimulateError(seed int64) error {
	sim := NewSimulator(seed)
	return sim.Inject("random")
}
