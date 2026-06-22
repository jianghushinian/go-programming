package cronx

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/robfig/cron/v3"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	"cronx/internal/pkg/store"
	"cronx/internal/watcher"
	_ "cronx/internal/watcher/all" // auto register all watchers
	"cronx/pkg/config"
	"cronx/pkg/distributedlock"
	"cronx/pkg/logx"
)

var (
	lockName          = "cronx:lock"
	jobStopTimeout    = 3 * time.Minute
	extendExpiration  = 5 * time.Second
	defaultExpiration = 10 * extendExpiration
)

// Config holds all necessary runtime configuration for the Cronx server.
type Config struct {
	MySQLConfig *config.MySQLConfig
	RedisConfig *config.RedisConfig
	Clientset   kubernetes.Interface
}

// CreateWatcherConfig builds the shared watcher.Config used by all watcher instances.
func (c *Config) CreateWatcherConfig() (*watcher.Config, error) {
	gormDB, err := config.NewMySQL(c.MySQLConfig)
	if err != nil {
		logx.Error(context.Background(), err.Error(), "Failed to create MySQL client")
		return nil, err
	}
	datastore := store.NewStore(gormDB)

	return &watcher.Config{Store: datastore, Clientset: c.Clientset}, nil
}

// NewCronx initializes the full Cronx runtime.
func (c *Config) NewCronx(ctx context.Context) (*Cronx, error) {
	rdb, err := config.NewRedis(c.RedisConfig)
	if err != nil {
		logx.Error(ctx, err.Error(), "Failed to create Redis client")
		return nil, err
	}

	// Configure cron scheduler.
	logger := newCronLogger()
	scheduler := cron.New(
		cron.WithSeconds(),
		cron.WithLogger(logger),
		cron.WithChain(cron.SkipIfStillRunning(logger), cron.Recover(logger)),
	)

	// Initialize distributed lock.
	pool := goredis.NewPool(rdb)
	lockOpts := []redsync.Option{
		redsync.WithRetryDelay(50 * time.Microsecond),
		redsync.WithTries(3),
		redsync.WithExpiry(defaultExpiration),
	}
	mutex := redsync.New(pool).NewMutex(lockName, lockOpts...)
	locker := distributedlock.New(mutex, extendExpiration)

	// Build watcher configuration.
	cfg, err := c.CreateWatcherConfig()
	if err != nil {
		return nil, err
	}

	// Build Cronx.
	cx := &Cronx{scheduler: scheduler, locker: locker, config: cfg}

	return cx, nil
}

// Cronx represents the main job scheduler engine.
type Cronx struct {
	scheduler *cron.Cron            // scheduler by cron
	locker    *distributedlock.Lock // distributed lock
	config    *watcher.Config
}

// addWatchers registers all discovered watchers into the scheduler.
// Each watcher may implement its own cron spec via ISpec; otherwise defaults to 10s.
func (c *Cronx) addWatchers(ctx context.Context) error {
	for name, w := range watcher.ListWatchers() {
		if err := w.Init(ctx, c.config); err != nil {
			logx.Error(ctx, err.Error(), "Failed to construct watcher", "watcher", name)
			return err
		}

		spec := watcher.Every10Seconds
		if obj, ok := w.(watcher.ISpec); ok {
			spec = obj.Spec()
		}

		if _, err := c.scheduler.AddJob(spec, w); err != nil {
			logx.Error(ctx, err.Error(), "Failed to add job to the cron", "watcher", name)
			return err
		}
	}

	return nil
}

// Run starts the cron scheduler and blocks until stopCh is closed.
// Only one Cronx instance can run at a time due to the distributed lock.
func (c *Cronx) Run(stopCh <-chan struct{}) error {
	ctx := wait.ContextForChannel(stopCh)

	// Attempt to acquire the distributed lock in a loop.
	ticker := time.NewTicker(defaultExpiration + (5 * time.Second))
	defer ticker.Stop()
	for {
		err := c.locker.Lock(ctx)
		if err == nil {
			logx.Info(ctx, "Successfully acquired lock", "lockName", lockName)
			break
		}

		select {
		case <-ctx.Done():
			logx.Info(ctx, "Shutting down acquired lock")
			return nil
		case <-ticker.C:
			logx.Debug(ctx, "Failed to acquire lock", "lockName", lockName, "err", err)
		}
	}

	// Registers all watchers.
	if err := c.addWatchers(ctx); err != nil {
		return err
	}
	logx.Info(ctx, "Successfully registers watchers")

	// Start scheduled jobs.
	c.scheduler.Start()
	logx.Info(ctx, "Successfully started cronx server")

	// Block until shutdown signal is received.
	<-stopCh

	logx.Info(ctx, "Stopping...")
	c.stop()

	return nil
}

// stop gracefully shuts down all cron jobs and releases the distributed lock.
func (c *Cronx) stop() {
	ctx := c.scheduler.Stop()

	// Wait for all jobs to stop, or time out.
	select {
	case <-ctx.Done():
	case <-time.After(jobStopTimeout):
		logx.Error(ctx, "Context was not done immediately", "timeout", jobStopTimeout.String())
	}

	// Attempt to release lock.
	if ok, err := c.locker.Unlock(); !ok || err != nil {
		logx.Error(ctx, "Failed to unlock", "err", err, "status", ok)
	}
}
