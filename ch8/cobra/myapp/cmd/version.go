/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v", "ver"},
	GroupID: "info",
	Short:   "Show version",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version called")
	},
	Args: cobra.MaximumNArgs(2), // 位置参数多于 2 个则报错
	// Args: func(cmd *cobra.Command, args []string) error {
	// 	if len(args) < 1 {
	// 		return errors.New("requires at least one arg")
	// 	}
	// 	if len(args) > 4 {
	// 		return errors.New("the number of args cannot exceed 4")
	// 	}
	// 	return nil
	// },
	PreRunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("version PreRunE called")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("version called")
		fmt.Printf("cfgFile: %s\n", cfgFile)
		return nil
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("version PostRun called")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
