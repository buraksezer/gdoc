package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gdoc",
	Short: "Search GoDoc.org via command-line",
}

func Execute() {
	rootCmd.Execute()
}
