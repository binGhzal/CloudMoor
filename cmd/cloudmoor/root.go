package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cloudmoor",
	Short: "CloudMoor - Remote Storage Mounting Platform",
	Long: `CloudMoor provides secure, unified access to cloud storage providers
through a persistent daemon and intuitive CLI.`,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
