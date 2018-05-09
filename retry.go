package retry

import (
	"context"
	"math/rand"
	"time"
)

// Do runs a function until the BackoffStrategy is exhausted or the function
// returns nil
func Do(backoff BackoffStrategy, funcToRetry func() error) (err error) {
	for {
		err := funcToRetry()
		if err == nil {
			return nil
		}
		nextDuration, toContinue := backoff()
		if !toContinue {
			return err
		}
		time.Sleep(nextDuration)
	}
}

// DoWithContext runs a function until the BackoffStrategy is exhausted, until
// the context is done, or until the function returns nil. Please note: it is
// the responsibility of the called function to ensure context is obeyed if it
// is to exit in a timely manner once the context is done.
func DoWithContext(ctx context.Context, backoff BackoffStrategy, funcToRetry func(ctx context.Context) error) (err error) {
	for {
		err := funcToRetry(ctx)
		if err == nil {
			return nil
		}
		nextDuration, toContinue := backoff()
		if !toContinue {
			return err
		}
		select {
		case <-time.NewTimer(nextDuration).C:
		case <-ctx.Done():
			// context cancelled, return the last error we got
			return err
		}
	}
}

// BackoffStrategy represents a function that returns successive wait durations
// and a bool representing whether or not to continue
type BackoffStrategy func() (time.Duration, bool)

// ConstantBackoff always returns the same duration until maxAttempts - 1 is
// reached
func ConstantBackoff(maxAttempts int, delay time.Duration) BackoffStrategy {
	attempts := 0
	return func() (time.Duration, bool) {
		attempts++
		return delay, attempts < maxAttempts-1
	}
}

// ExponentialBackoff implements the exponential backoff algorithm to gradually
// increase backoff time with jitter until maxAttempts -1 is reached.
//
// https://en.wikipedia.org/wiki/Exponential_backoff
// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
func ExponentialBackoff(maxAttempts int, initialDelay time.Duration, maxDelay time.Duration) BackoffStrategy {
	if initialDelay <= 0 {
		panic("ExponentialBackoff requires positive duration")
	}
	attempts := 0
	return func() (time.Duration, bool) {
		attempts++
		return min(
			maxDelay,
			time.Duration(rand.Int63n((1 << uint((attempts - 1)) * int64(initialDelay)))),
		), attempts < maxAttempts-1
	}
}

func min(a, b time.Duration) time.Duration {
	if a > b {
		return b
	}
	return a
}

func init() {
	rand.Seed(time.Now().UnixNano())
}