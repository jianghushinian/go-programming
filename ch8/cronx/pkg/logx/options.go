package logx

import "github.com/spf13/pflag"

// Options defines options for logx.
type Options struct {
	Level       string   `json:"level" yaml:"level" mapstructure:"level"`    // debug, info, warn, error
	Format      string   `json:"format" yaml:"format" mapstructure:"format"` // json, text
	OutputPaths []string `json:"output" yaml:"output" mapstructure:"output"` // default stdout
}

func NewOptions() *Options {
	return &Options{
		Level:       "info",
		Format:      "text",
		OutputPaths: []string{"stdout"},
	}
}

// Validate verifies flags passed to Options.
func (o *Options) Validate() []error {
	var errs []error
	return errs
}

// AddFlags adds flags related to logx for a specific server to the specified FlagSet.
func (o *Options) AddFlags(fs *pflag.FlagSet, prefixes ...string) {
	fs.StringVar(&o.Level, "log.level", "info", "LogOptions level: debug|info|warn|error")
	fs.StringVar(&o.Format, "log.format", "text", "LogOptions format: text|json")
	// e.g., --log.output=stdout,/var/log/cronx.log
	fs.StringSliceVar(&o.OutputPaths, "log.output", []string{"stdout"}, "LogOptions output: stdout or filepath")
}
