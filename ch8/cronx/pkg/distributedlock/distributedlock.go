package distributedlock

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"

	"cronx/pkg/logx"
)

// Lock is a wrapper around redsync.Mutex to manage distributed locks.
type Lock struct {
	mutex      *redsync.Mutex // The underlying redsync mutex for lock management
	stopCh     chan struct{}  // Channel to signal when to stop lock refresh
	extendTick time.Duration  // Duration between lock refresh attempts
}

// New creates and returns a new distributed Lock instance.
func New(mutex *redsync.Mutex, extendTick time.Duration) *Lock {
	return &Lock{
		mutex:      mutex,
		stopCh:     make(chan struct{}),
		extendTick: extendTick,
	}
}

// Lock acquires the lock and starts automatic lock refresh.
func (l *Lock) Lock(ctx context.Context) error {
	if err := l.mutex.LockContext(ctx); err != nil {
		return err
	}
	// Start the auto-refresh in a new goroutine
	go l.autoRefresh(ctx)
	return nil
}

// TryLock attempts to acquire the lock without blocking and starts automatic refresh if successful.
func (l *Lock) TryLock(ctx context.Context) error {
	err := l.mutex.TryLockContext(ctx)
	if err != nil {
		return err
	}
	go l.autoRefresh(ctx)
	return nil
}

// Unlock stops the lock refresh process and releases the lock.
func (l *Lock) Unlock() (bool, error) {
	close(l.stopCh)
	return l.mutex.Unlock()
}

// RunWithLock executes the provided function while holding the lock.
func (l *Lock) RunWithLock(ctx context.Context, fn func(ctx context.Context)) error {
	if err := l.Lock(ctx); err != nil {
		return err
	}
	defer l.Unlock()

	fn(ctx)
	return nil
}

// autoRefresh keeps extending the lock expiration at regular intervals.
func (l *Lock) autoRefresh(ctx context.Context) {
	ticker := time.NewTicker(l.extendTick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if ok, err := l.mutex.ExtendContext(ctx); !ok || err != nil {
				logx.Error(ctx, "Failed to extend lock", "err", err, "status", ok)
			} else {
				logx.Info(ctx, "Successes to extend lock")
			}
		case <-l.stopCh:
			return
		}
	}
}
