package cmd

import (
	"bytes"
	"fmt"
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
	readURL      = "https://godoc.org"
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
	defer func() {
		if derr := db.Close(); derr != nil {
			logger.Println("failed to close boltdb: ", derr)
		}
	}()

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
	if err := getDocument(pkgpath); err != nil {
		logger.Fatal(fmt.Sprintf("failed to get document: %v", err))
	}
}

func getDocument(pkgpath string) error {
	target, _ := url.Parse(readURL)
	target.Path = path.Join("/", pkgpath)
	req, err := http.NewRequest(http.MethodGet, target.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/plain")
	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if berr := r.Body.Close(); err != nil {
			logger.Println("failed to close request body:", berr)
		}
	}()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		msg := strings.TrimSpace(string(data))
		return fmt.Errorf("godoc.org returned HTTP %d: %s", r.StatusCode, msg)
	}

	if disablePager {
		fmt.Println(string(data))
		return nil
	}
	// GDOC_PAGER is useful for setting special pager command for gdoc.
	cmd := os.Getenv("GDOC_PAGER")
	if cmd == "" {
		cmd = os.Getenv("PAGER")
	}
	if cmd == "" {
		logger.Fatal("no available pager found")
	}
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stdin = bytes.NewReader(data)
	return c.Run()
}
