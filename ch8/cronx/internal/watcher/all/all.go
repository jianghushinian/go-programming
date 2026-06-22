package all

import (
	// Importing all watchers solely for the side effects of their init() functions,
	// ensuring that each watcher self-registers into the global registry.
	_ "cronx/internal/watcher/job"
)
