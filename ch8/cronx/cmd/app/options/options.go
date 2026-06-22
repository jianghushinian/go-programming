package options

import (
	"context"
	"os"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"cronx/internal/cronx"
	"cronx/pkg/config"
	"cronx/pkg/logx"
	genericoptions "cronx/pkg/options"
	kubeutil "cronx/pkg/util/kube"
)

const (
	// UserAgent is used when sending requests to the Kubernetes API server.
	UserAgent = "cronx"
)

// Options contains everything necessary to create and run a cronx server.
type Options struct {
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
	RedisOptions *genericoptions.RedisOptions `json:"redis" mapstructure:"redis"`
	// Path to kubeconfig file with authorization and master location information.
	Kubeconfig string        `json:"kubeconfig" mapstructure:"kubeconfig"`
	LogOptions *logx.Options `json:"log" mapstructure:"log"`
}

// NewOptions returns an Options instance populated with default values.
func NewOptions() *Options {
	return &Options{
		MySQLOptions: genericoptions.NewMySQLOptions(),
		RedisOptions: genericoptions.NewRedisOptions(),
		LogOptions:   logx.NewOptions(),
	}
}

// AddFlags binds the options to CLI flags.
func (o *Options) AddFlags(fs *pflag.FlagSet) {
	o.MySQLOptions.AddFlags(fs)
	o.RedisOptions.AddFlags(fs)
	o.LogOptions.AddFlags(fs)

	// Path to kubeconfig for job watcher.
	fs.StringVar(&o.Kubeconfig, "kubeconfig", defaultKubeconfigPath(), "Path to kubeconfig file with authorization and master location information.")
}

// Validate validates all the required options.
func (o *Options) Validate() error {
	var errs []error

	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)

	return utilerrors.NewAggregate(errs)
}

// ApplyTo fills up cronx config with options.
func (o *Options) ApplyTo(c *cronx.Config) error {
	if err := copier.Copy(c.MySQLConfig, o.MySQLOptions); err != nil {
		return err
	}
	if err := copier.Copy(c.RedisConfig, o.RedisOptions); err != nil {
		return err
	}
	return nil
}

// Config return a cronx config object.
func (o *Options) Config() (*cronx.Config, error) {
	// First try in-cluster configuration (when running inside Kubernetes)
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig file for local development or external execution
		logx.Warn(context.Background(), "Failed to create in-cluster config", "err", err.Error())
		kubeconfig, err = clientcmd.BuildConfigFromFlags("", o.Kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	kubeutil.SetDefaultClientOptions(kubeutil.AddUserAgent(kubeconfig, UserAgent))

	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	c := &cronx.Config{
		Clientset:   clientset,
		MySQLConfig: &config.MySQLConfig{},
		RedisConfig: &config.RedisConfig{},
	}

	if err := o.ApplyTo(c); err != nil {
		return nil, err
	}

	return c, nil
}

// defaultKubeconfigPath determines the default path to the kubeconfig file.
func defaultKubeconfigPath() string {
	// Find home directory.
	homeDir, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return filepath.Join(homeDir, ".kube", "config")
}
