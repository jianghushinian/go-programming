/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// gendocsCmd represents the gendocs command
var genDocsCmd = &cobra.Command{
	Use:   "gendocs",
	Short: "Generate CLI documentation in Markdown format",
	RunE: func(cmd *cobra.Command, args []string) error {
		outputDir := "./docs"
		// 若不存在则自动创建
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}
		// 为整个命令树生成 Markdown 文档
		if err := doc.GenMarkdownTree(rootCmd, outputDir); err != nil {
			return err
		}
		fmt.Println("Markdown docs generated at:", outputDir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(genDocsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gendocsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gendocsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
