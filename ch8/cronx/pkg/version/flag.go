package version

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

// AddFlags adds a --version flag to the command.
func AddFlags(fs *pflag.FlagSet) {
	fs.Bool("version", false, "Print version information and exit.")
}

// CheckVersionFlag checks if user passed --version and exits.
func CheckVersionFlag(fs *pflag.FlagSet) {
	showVersion, _ := fs.GetBool("version")
	if showVersion {
		fmt.Println(String())
		os.Exit(0)
	}
}
