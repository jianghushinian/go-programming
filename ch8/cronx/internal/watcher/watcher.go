package watcher

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"

	reflectutil "cronx/pkg/util/reflect"
)

const (
	// Every10Seconds defines a common cron expression used by internal watchers.
	// It schedules the job to run once every 10 seconds.
	Every10Seconds = "@every 10s"
)

// Watcher represents a pluggable cron-based watcher component.
// Each watcher is responsible for its own initialization and task execution.
// It embeds cron.Job so that it can be executed directly by the cron scheduler.
type Watcher interface {
	// Init initializes the watcher with the provided configuration.
	// This is typically where dependencies are prepared or caches are loaded.
	Init(ctx context.Context, config *Config) error

	// Job Run method comes from cron.Job interface.
	cron.Job
}

// ISpec provides an optional method for watchers to define their own cron schedule.
// Implementing this interface allows a watcher to self-declare its own schedule.
type ISpec interface {
	// Spec return the spec for a cron job.
	// There are two cron spec formats in common usage:
	// - standard cron format: https://en.wikipedia.org/wiki/Cron
	// - quartz scheduler format: http://www.quartz-scheduler.org/documentation/quartz-2.3.0/tutorials/tutorial-lesson-06.html
	// This method is optional for a watcher.
	Spec() string
}

var (
	// registryLock protects access to the watcher registry for concurrent read/write.
	registryLock = new(sync.RWMutex)

	// registry stores all registered watchers keyed by their struct name.
	// Watchers are typically registered via init() in their respective packages.
	registry = make(map[string]Watcher)
)

// Register adds a watcher into the global registry.
// It panics if a watcher with the same struct name has already been registered.
// Typically called inside an init() function of the watcher implementation.
func Register(watcher Watcher) {
	registryLock.Lock()
	defer registryLock.Unlock()

	name := reflectutil.StructName(watcher)
	if _, ok := registry[name]; ok {
		panic("duplicate watcher entry: " + name)
	}

	registry[name] = watcher
}

// ListWatchers returns a snapshot of all registered watchers.
// Callers can iterate through the returned map to initialize or schedule them.
func ListWatchers() map[string]Watcher {
	registryLock.RLock()
	defer registryLock.RUnlock()

	return registry
}
