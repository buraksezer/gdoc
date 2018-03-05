package cmd

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	releaseVersion = "0.1"
	dbname         = ".gdoc.db"
)

var rootCmd = &cobra.Command{
	Use:     "gdoc",
	Short:   "Search GoDoc.org via command-line",
	Version: fmt.Sprintf("%s with %s", releaseVersion, runtime.Version()),
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
