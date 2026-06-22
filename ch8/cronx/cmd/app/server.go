package app

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	genericapiserver "k8s.io/apiserver/pkg/server"

	"cronx/cmd/app/options"
	"cronx/pkg/logx"
	"cronx/pkg/version"
)

// configFile stores the path of the configuration file specified via CLI.
var configFile string

// NewCronxCommand creates and returns the root `cronx` command.
// This command initializes configuration, binds flags, and launches the server.
func NewCronxCommand() *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:   "cronx",
		Short: "Launch a cron-like asynchronous job processing server",
		Long: `The cronx server is responsible for executing some async tasks 
like linux cronjob. You can add Cron(github.com/robfig/cron) jobs on the given schedule
use the Cron spec format.`,
		SilenceUsage: true,

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			version.CheckVersionFlag(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
	}

	// Initialize configuration when the command starts.
	cobra.OnInitialize(initConfig(&configFile, "CRONX", searchDirs(), defaultConfigName))
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the cronx configuration file.")

	// Bind option flags (log, storage, etc).
	opts.AddFlags(cmd.PersistentFlags())

	// Add --version flag.
	version.AddFlags(cmd.PersistentFlags())

	return cmd
}

// run starts the cronx server. This function should block and not exit
// until the server shutdown due to a system signal.
func run(opts *options.Options) error {
	// SetupSignalContext returns a context canceled when SIGTERM/SIGINT is received.
	ctx := genericapiserver.SetupSignalContext()

	// Load configuration values from viper into opts.
	if err := viper.Unmarshal(opts); err != nil {
		return fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	// Initialize global logger using logx.
	logx.Init(opts.LogOptions)
	defer logx.Sync()

	// Convert options into a full configuration object.
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	// Create the Cronx server instance.
	cx, err := cfg.NewCronx(ctx)
	if err != nil {
		return err
	}

	logx.Info(ctx, "Starting cronx...")

	// Start the server and block until context is canceled.
	return cx.Run(ctx.Done())
}
