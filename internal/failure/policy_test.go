package failure

import (
	"errors"
	"testing"
)

func TestRetryWithPolicySuccess(t *testing.T) {
	policy := NewRetryPolicy(3, 1)
	attempts := 0
	err := RetryWithPolicy(&policy, func() error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestSimulatorDeterministic(t *testing.T) {
	a := NewSimulator(42).Inject("random")
	b := NewSimulator(42).Inject("random")
	if (a == nil) != (b == nil) {
		t.Fatalf("expected same outcome for same seed")
	}
}

func TestSimulatorKinds(t *testing.T) {
	if NewSimulator(1).Inject("timeout") != ErrTimeout {
		t.Fatal("expected timeout")
	}
	if NewSimulator(1).Inject("rate_limit") != ErrRateLimit {
		t.Fatal("expected rate limit")
	}
}
