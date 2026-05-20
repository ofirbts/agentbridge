package failure

import (
	"errors"
	"math/rand"
)

var (
	ErrTimeout     = errors.New("simulated timeout")
	ErrRateLimit   = errors.New("simulated rate limit 429")
	ErrServerError = errors.New("simulated server error 500")
)

type Simulator struct {
	rng *rand.Rand
}

func NewSimulator(seed int64) *Simulator {
	if seed == 0 {
		seed = 42
	}
	return &Simulator{rng: rand.New(rand.NewSource(seed))}
}

func (s *Simulator) Inject(kind string) error {
	switch kind {
	case "timeout":
		return ErrTimeout
	case "rate_limit":
		return ErrRateLimit
	case "server_error":
		return ErrServerError
	case "random":
		if s.rng.Intn(2) == 0 {
			return ErrSimulated
		}
		return nil
	default:
		return ErrSimulated
	}
}
