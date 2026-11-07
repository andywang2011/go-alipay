package alipay

import (
	"context"
	"errors"
	"io"
	"sync"
	"testing"
	"time"
)

func TestDefaultRetryConfig(t *testing.T) {
	c := DefaultRetryConfig()
	if c.MaxAttempts != 2 || c.InitialDelay <= 0 {
		t.Fatal(c)
	}
}

func TestRetryManager_Success(t *testing.T) {
	rm := newRetryManager(DefaultRetryConfig())
	n := 0
	err := rm.execute(context.Background(), func(context.Context) error {
		n++
		return nil
	})
	if err != nil || n != 1 {
		t.Fatalf("n=%d err=%v", n, err)
	}
}

func TestRetryManager_RetryThenSuccess(t *testing.T) {
	rm := newRetryManager(RetryConfig{
		MaxAttempts:  3,
		InitialDelay: time.Millisecond,
		MaxDelay:     time.Millisecond,
		Multiplier:   2,
	})
	rm.sleep = func(ctx context.Context, d time.Duration) error { return nil }
	n := 0
	err := rm.execute(context.Background(), func(context.Context) error {
		n++
		if n < 2 {
			return io.EOF
		}
		return nil
	})
	if err != nil || n != 2 {
		t.Fatalf("n=%d err=%v", n, err)
	}
}

func TestRetryManager_NonRetryableStops(t *testing.T) {
	rm := newRetryManager(RetryConfig{MaxAttempts: 3, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, Multiplier: 1})
	n := 0
	err := rm.execute(context.Background(), func(context.Context) error {
		n++
		return NewAlipayErrorFromBase("40004", "fail", "", "")
	})
	if err == nil || n != 1 {
		t.Fatalf("n=%d err=%v", n, err)
	}
}

func TestRetryManager_ContextCanceled(t *testing.T) {
	rm := newRetryManager(DefaultRetryConfig())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := rm.execute(ctx, func(context.Context) error { return io.EOF })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
}

func TestIsRetryableHTTPError(t *testing.T) {
	if !isRetryableHTTPError(io.EOF) {
		t.Fatal("EOF should be retryable")
	}
	if isRetryableHTTPError(NewAlipayErrorFromBase("1", "m", "", "")) {
		t.Fatal("AlipayError should not retry")
	}
	if !isRetryableHTTPError(errors.New("alipay http status 503: boom")) {
		t.Fatal("503 should retry")
	}
	if isRetryableHTTPError(errors.New("alipay http status 400: bad")) {
		t.Fatal("400 should not retry")
	}
}

func TestNormalizeRetry(t *testing.T) {
	r := RetryConfig{}
	normalizeRetry(&r)
	if r.MaxAttempts != DefaultRetryConfig().MaxAttempts {
		t.Fatal(r)
	}
}

func TestConfig_Validate_AppliesRetry(t *testing.T) {
	cfg := testConfig()
	cfg.RetryConfig = RetryConfig{}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	if cfg.RetryConfig.MaxAttempts != DefaultRetryConfig().MaxAttempts {
		t.Fatalf("retry=%+v", cfg.RetryConfig)
	}
}

func TestSleepContext(t *testing.T) {
	if err := sleepContext(context.Background(), 0); err != nil {
		t.Fatalf("zero delay: %v", err)
	}
	if err := sleepContext(context.Background(), time.Millisecond); err != nil {
		t.Fatalf("elapsed: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := sleepContext(ctx, time.Hour); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled: %v", err)
	}
	if err := sleepContext(ctx, 0); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled zero delay: %v", err)
	}
}

func TestRetryManager_JitterInRangeAndConcurrent(t *testing.T) {
	rm := newRetryManager(DefaultRetryConfig())
	if rm.jitter(0) != 0 {
		t.Fatal("zero delay must give zero jitter")
	}
	d := 100 * time.Millisecond
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				if got := rm.jitter(d); got < d*8/10 || got >= d*12/10 {
					t.Errorf("jitter %v outside ±20%% of %v", got, d)
					return
				}
			}
		}()
	}
	wg.Wait()
	if got := rm.jitter(3); got < 2 || got > 4 {
		t.Fatalf("tiny delay jitter = %v", got)
	}
}
