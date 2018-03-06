package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	releaseVersion = "0.1"
	dbname         = ".gdoc.db"
)

var logger = log.New(os.Stderr, "error: ", 0)

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

// Execute uses the args (os.Args[1:] by default) and run through the command tree finding
// appropriate matches for commands and then corresponding flags.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Println("failed to run gdoc:", err)
	}
}
