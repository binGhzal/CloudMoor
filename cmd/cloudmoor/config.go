package main

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CloudMoor configuration",
	Long:  "Commands for managing CloudMoor configuration, including vault operations.",
}

func init() {
	configCmd.AddCommand(vaultCmd)
}
