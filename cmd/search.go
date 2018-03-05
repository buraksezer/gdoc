package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type Result struct {
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	ImportCount int     `json:"import_count"`
	Synopsis    string  `json:"synopsis"`
	Stars       int     `json:"stars"`
	Score       float64 `json:"score"`
}

type Results struct {
	Results []Result
}

var (
	interactive    bool
	maxResultCount int
	client         = &http.Client{}
	searchUrl      = "https://api.godoc.org/search"
)

var searchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "Search given keyword on GoDoc.org",
	Args:  cobra.MinimumNArgs(1),
	Run:   search,
}

func init() {
	searchCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "enable interactive mode")
	searchCmd.Flags().IntVarP(&maxResultCount, "count", "c", 10, "sets maximum result count")
	rootCmd.AddCommand(searchCmd)
}

func search(_ *cobra.Command, args []string) {
	target, _ := url.Parse(searchUrl)
	target.RawQuery = "q=" + args[0]
	r, err := client.Get(target.String())
	if err != nil {
		logger.Fatalf("failed to search: %v", err)
	}
	defer r.Body.Close()

	results := &Results{}
	if err := json.NewDecoder(r.Body).Decode(results); err != nil {
		logger.Fatalf("failed to search: %v", err)
	}

	if len(results.Results) == 0 {
		logger.Fatalf("nothing found")
	}

	pkgs := make(map[int]string)
	for index, res := range results.Results {
		if interactive {
			fmt.Printf("==> (%d) %s\n", index+1, res.Path)
			pkgs[index+1] = res.Path
		} else {
			fmt.Printf("==> %s\n", res.Path)
		}
		fmt.Printf("==> imports: %d stars: %d\n", res.ImportCount, res.Stars)
		if len(res.Synopsis) > 0 {
			fmt.Printf("%s\n", res.Synopsis)
		}
		if maxResultCount != 0 && index+1 >= maxResultCount {
			break
		}
		fmt.Printf("\n")
	}

	if interactive {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("\nGive a number to read the document: \n")
		snum, err := reader.ReadString('\n')
		if err != nil {
			logger.Fatal(err)
		}
		snum = strings.TrimSpace(snum)
		num, err := strconv.Atoi(snum)
		if err != nil {
			logger.Fatal(err)
		}
		pkg, ok := pkgs[num]
		if !ok {
			logger.Fatal("invalid index")
		}
		read(nil, []string{pkg})
	}
}
