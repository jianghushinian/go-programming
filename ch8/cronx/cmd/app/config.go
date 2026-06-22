package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// defaultHomeDir defines the default directory under $HOME
	// where cronx configuration files are stored.
	defaultHomeDir = ".cronx"

	// defaultConfigName specifies the default name of the cronx
	// configuration file when no file is explicitly provided.
	defaultConfigName = "cronx.yaml"
)

// initConfig returns a function used by Cobra's initialization hook.
// It configures Viper to load settings from a file, environment variables,
// and predefined search directories.
//
// Parameters:
//   - configFile:     optional path to a specific config file
//   - envPrefix:      prefix used when reading environment variables
//   - loadDirs:       directories where Viper will search for config files
//   - defaultConfigName: default filename used when searching for configs
func initConfig(configFile *string, envPrefix string, loadDirs []string, defaultConfigName string) func() {
	return func() {
		if configFile != nil {
			// Use configuration file provided via CLI flag.
			viper.SetConfigFile(*configFile)
		} else {
			// Otherwise, search for the configuration file in the given directories.
			// Viper looks for files named <defaultConfigName>.yaml, .yml, or .json.
			for _, dir := range loadDirs {
				viper.AddConfigPath(dir)
			}
			viper.SetConfigType("yaml")
			viper.SetConfigName(defaultConfigName)
		}

		// Search for any environment variable prefixed and load them in.
		viper.AutomaticEnv()
		// Tells Viper to use this prefix when reading environment variables.
		viper.SetEnvPrefix(envPrefix)

		// Convert "." and "-" to "_" when matching environment variables.
		// e.g., log.level -> LOG_LEVEL
		replacer := strings.NewReplacer(".", "_", "-", "_")
		viper.SetEnvKeyReplacer(replacer)

		// Find and read the config file.
		_ = viper.ReadInConfig()
	}
}

func searchDirs() []string {
	// Find home directory.
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return []string{filepath.Join(homeDir, defaultHomeDir), "."}
}

func filePath() string {
	// Find home directory.
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return filepath.Join(homeDir, defaultHomeDir, defaultConfigName)
}
