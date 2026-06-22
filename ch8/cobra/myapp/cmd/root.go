/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello Cobra!")
		toggle, _ := cmd.Flags().GetBool("toggle")
		fmt.Printf("toggle: %t\n", toggle)
		fmt.Printf("cfgFile: %s\n", cfgFile)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig) // 与 Viper 集成

	// rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("==== MyApp Help ====")
	// 	fmt.Println("Usage:")
	// 	cmd.Println(cmd.UseLine())
	// 	fmt.Println("\nAvailable Commands:")
	// 	for _, command := range cmd.Commands() {
	// 		fmt.Printf("\t%s\n", command.Use)
	// 	}
	// })

	// 	rootCmd.SetHelpTemplate(`
	// Custom Help for {{.UseLine}}
	//
	// Description:
	//   {{.Short}}
	//
	// Flags:
	// {{.Flags.FlagUsages}}
	// `)

	// customHelpCmd := &cobra.Command{
	// 	Use:   "help",
	// 	Short: "My custom help",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		fmt.Println("This is custom help output.")
	// 	},
	// }
	//
	// rootCmd.SetHelpCommand(customHelpCmd)

	// rootCmd.SetUsageFunc(func(cmd *cobra.Command) error {
	// 	fmt.Printf("Custom Usage for %s:\n", cmd.Name())
	// 	fmt.Println("Example:")
	// 	fmt.Println("  myapp version")
	// 	return nil
	// })

	rootCmd.SetUsageTemplate(`Usage: myapp serve [OPTIONS]
Start the server with custom configs.
`)

	rootCmd.AddGroup(&cobra.Group{
		ID:    "info",
		Title: "Information Commands",
	})

	rootCmd.AddGroup(&cobra.Group{
		ID:    "ops",
		Title: "Operations",
	})

	rootCmd.SuggestionsMinimumDistance = 1
	// rootCmd.DisableSuggestions = true

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
