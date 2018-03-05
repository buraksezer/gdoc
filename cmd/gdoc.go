package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
)

const dbname = ".gdoc.db"

var rootCmd = &cobra.Command{
	Use:   "gdoc",
	Short: "Search GoDoc.org via command-line",
}

func findGdocDB() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("$HOME is empty")
	}
	return path.Join(home, dbname), nil
}

func Execute() {
	rootCmd.Execute()
}
