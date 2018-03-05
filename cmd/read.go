package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	bolt "github.com/coreos/bbolt"
	"github.com/spf13/cobra"
)

var (
	readUrl      = "https://godoc.org"
	disablePager bool
	pkgAlias     bool
)

var readCmd = &cobra.Command{
	Use:   "read [pkgpath]",
	Short: "Get the documentation for a package",
	Args:  cobra.MinimumNArgs(1),
	Run:   read,
}

func init() {
	readCmd.Flags().BoolVar(&disablePager, "disable-pager", false, "disables piping package documentation to the pager")
	readCmd.Flags().BoolVarP(&pkgAlias, "alias", "a", false, "retrieve package documentation by its alias")
	rootCmd.AddCommand(readCmd)
}

func getPkgpath(name string) (string, error) {
	dpath, err := findGdocDB()
	if err != nil {
		return "", err
	}
	db, err := bolt.Open(dpath, 0600, nil)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var pkgpath string
	res := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("alias"))
		if b == nil {
			return fmt.Errorf("alias could not be found")
		}
		v := b.Get([]byte(name))
		if len(v) == 0 {
			return fmt.Errorf("alias could not be found")
		}
		pkgpath = string(v)
		return nil
	})
	return pkgpath, res
}

func read(_ *cobra.Command, args []string) {
	var pkgpath string
	var err error
	if !pkgAlias {
		pkgpath = args[0]
	} else {
		pkgpath, err = getPkgpath(args[0])
		if err != nil {
			log.Fatal(err)
		}
	}
	target, _ := url.Parse(readUrl)
	target.Path = path.Join("/", pkgpath)
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
