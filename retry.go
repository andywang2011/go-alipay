package alipay

import (
	"context"
	"errors"
	"io"
	"math/rand/v2"
	"net"
	"strings"
	"time"
)

// RetryConfig controls exponential-backoff retry for transient HTTP failures.
type RetryConfig struct {
	// MaxAttempts is the total number of attempts including the first.
	MaxAttempts int
	// InitialDelay is the base delay before the first retry.
	InitialDelay time.Duration
	// MaxDelay caps exponential backoff delay.
	MaxDelay time.Duration
	// Multiplier scales delay after each failed attempt.
	Multiplier float64
}

// DefaultRetryConfig returns HTTP retry defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  2,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     2 * time.Second,
		Multiplier:   2.0,
	}
}

type retryManager struct {
	config RetryConfig
	// int64n returns a value in [0, n); defaults to math/rand/v2 (goroutine-safe).
	int64n func(n int64) int64
	sleep  func(ctx context.Context, d time.Duration) error
}

func newRetryManager(cfg RetryConfig) *retryManager {
	return &retryManager{
		config: cfg,
		int64n: rand.Int64N,
		sleep:  sleepContext,
	}
}

// sleepContext waits for d or until ctx is done, releasing the timer either way.
func sleepContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func (rm *retryManager) execute(ctx context.Context, op func(ctx context.Context) error) error {
	var lastErr error
	delay := rm.config.InitialDelay
	for attempt := 0; attempt < rm.config.MaxAttempts; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := op(ctx); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if attempt >= rm.config.MaxAttempts-1 || !isRetryableHTTPError(lastErr) {
			return lastErr
		}
		sleep := rm.jitter(delay)
		if err := rm.sleep(ctx, sleep); err != nil {
			return err
		}
		delay = time.Duration(float64(delay) * rm.config.Multiplier)
		if delay > rm.config.MaxDelay {
			delay = rm.config.MaxDelay
		}
	}
	return lastErr
}

func (rm *retryManager) jitter(delay time.Duration) time.Duration {
	if delay <= 0 {
		return 0
	}
	jitterRange := delay / 5
	if jitterRange == 0 {
		jitterRange = 1
	}
	j := time.Duration(rm.int64n(int64(jitterRange * 2)))
	return delay - jitterRange + j
}

func normalizeRetry(r *RetryConfig) {
	if r.MaxAttempts <= 0 {
		*r = DefaultRetryConfig()
	}
	if r.InitialDelay <= 0 {
		r.InitialDelay = DefaultRetryConfig().InitialDelay
	}
	if r.MaxDelay <= 0 {
		r.MaxDelay = DefaultRetryConfig().MaxDelay
	}
	if r.Multiplier <= 0 {
		r.Multiplier = DefaultRetryConfig().Multiplier
	}
	if r.MaxDelay < r.InitialDelay {
		r.MaxDelay = r.InitialDelay
	}
}

func isRetryableHTTPError(err error) bool {
	if err == nil {
		return false
	}
	if _, ok := IsAlipayError(err); ok {
		return false
	}
	if errors.Is(err, io.EOF) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	msg := err.Error()
	if strings.HasPrefix(msg, "alipay http status ") {
		var code int
		if _, scanErr := parseHTTPStatusPrefix(msg, &code); scanErr == nil && code >= 500 {
			return true
		}
		return false
	}
	if strings.Contains(msg, "alipay unmarshal:") {
		return false
	}
	if strings.Contains(msg, "alipay verify sign:") {
		return false
	}
	return false
}

func parseHTTPStatusPrefix(msg string, code *int) (int, error) {
	const prefix = "alipay http status "
	if !strings.HasPrefix(msg, prefix) {
		return 0, errors.New("not http status error")
	}
	rest := msg[len(prefix):]
	n := 0
	for n < len(rest) && rest[n] >= '0' && rest[n] <= '9' {
		*code = *code*10 + int(rest[n]-'0')
		n++
	}
	if n == 0 {
		return 0, errors.New("no status code")
	}
	return n, nil
}
