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
	Fork        bool    `json:"fork"`
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
	Use:              "search",
	Short:            "Search given keyword on GoDoc.org",
	TraverseChildren: true,
	Run:              search,
}

func init() {
	searchCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "enable interactive mode")
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

	pkgs := make(map[int]string)
	for index, res := range results.Results {
		if interactive {
			fmt.Printf("==> (%d) %s\n", index+1, res.Path)
			pkgs[index+1] = res.Path
		} else {
			fmt.Printf("==> %s\n", res.Path)
		}
		fmt.Printf("==> imports: %d stars: %d fork: %v\n", res.ImportCount, res.Stars, res.Fork)
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
			fmt.Println("error:", err)
			os.Exit(1)
		}
		snum = strings.TrimSpace(snum)
		num, err := strconv.Atoi(snum)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}
		pkg, ok := pkgs[num]
		if !ok {
			fmt.Println("error: invalid index")
			os.Exit(1)
		}
		fmt.Println("Getting documentation for", pkg)
		read(nil, []string{pkg})
	}
}
