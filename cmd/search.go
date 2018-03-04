package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

type Result struct {
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	ImportCount int     `json:"import_count"`
	Synopsis    string  `json:"synopsis"`
	Stars       int     `json:"stars"`
	Score       float64 `json:"score"`
	Fork        bool    `json:"fork"`
}

type Results struct {
	Results []Result
}

var (
	maxResultCount int
	client         = &http.Client{}
	searchUrl      = "https://api.godoc.org/search"
)

var searchCmd = &cobra.Command{
	Use:              "search",
	Short:            "Search given keyword on GoDoc.org",
	TraverseChildren: true,
	Run:              search,
}

func init() {
	searchCmd.Flags().IntVarP(&maxResultCount, "count", "c", 0, "sets maximum result count")
	rootCmd.AddCommand(searchCmd)
}

func search(_ *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("error: give a keyword to search")
		os.Exit(1)
	}
	target, _ := url.Parse(searchUrl)
	target.RawQuery = "q=" + args[0]
	r, err := client.Get(target.String())
	if err != nil {
		fmt.Println("error: failed to search:", err)
		os.Exit(1)
	}
	defer r.Body.Close()

	results := &Results{}
	if err := json.NewDecoder(r.Body).Decode(results); err != nil {
		fmt.Println("error: failed to search:", err)
		os.Exit(1)
	}
	for index, res := range results.Results {
		fmt.Printf("==> %s\n", res.Path)
		fmt.Printf("==> imports: %d stars: %d fork: %v\n", res.ImportCount, res.Stars, res.Fork)
		if len(res.Synopsis) > 0 {
			fmt.Printf("%s\n", res.Synopsis)
		}
		if maxResultCount != 0 && index+1 >= maxResultCount {
			return
		}
		fmt.Printf("\n")
	}
}
