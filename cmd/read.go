package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var (
	readUrl      = "https://godoc.org"
	disablePager bool
)

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Get the documentation for a package",
	Run:   read,
}

func init() {
	readCmd.Flags().BoolVar(&disablePager, "disable-pager", false, "disables piping package documentation to the pager")
	rootCmd.AddCommand(readCmd)
}

func read(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("error: give a keyword to search")
		os.Exit(1)
	}
	target, _ := url.Parse(readUrl)
	target.Path = path.Join("/", args[0])
	req, err := http.NewRequest(http.MethodGet, target.String(), nil)
	if err != nil {
		fmt.Println("Failed to get document:", err)
		os.Exit(1)
	}
	req.Header.Set("Accept", "text/plain")
	r, err := client.Do(req)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("error: failed to read response body: %v\n", err)
		os.Exit(1)
	}
	if r.StatusCode != http.StatusOK {
		fmt.Printf("error: godoc.org returned HTTP %d: %s\n",
			r.StatusCode, strings.TrimSpace(string(data)))
		os.Exit(1)
	}

	if disablePager {
		fmt.Println(string(data))
		return
	}

	c := exec.Command("/usr/bin/less")
	wr, err := c.StdinPipe()
	if err != nil {
		fmt.Printf("error: failed to get pipe: %v\n", err)
		os.Exit(1)
	}
	c.Stdout = os.Stdout
	reader := bytes.NewReader(data)
	go func() {
		_, cerr := io.Copy(wr, reader)
		if cerr != nil {
			fmt.Printf("error: failed to copy Stdin: %s\n", err)
		}
	}()
	if err := c.Run(); err != nil {
		fmt.Printf("error: failed to run less: %v\n", err)
		os.Exit(1)
	}
}
