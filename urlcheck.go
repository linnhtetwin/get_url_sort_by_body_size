package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/urfave/cli"
)

type result struct {
	url   string
	size  int
	error error
}

func main() {
	app := cli.NewApp()
	app.Name = "url-size"
	app.Usage = "visit a list of URLs and print their size"
	app.Action = func(c *cli.Context) error {
		urls := c.Args()
		results := make(chan result)

		for _, url := range urls {
			go func(url string) {
				res, err := http.Get(url)
				if err != nil {
					results <- result{url: url, error: err}
					return
				}
				defer res.Body.Close()
				results <- result{url: url, size: res.ContentLength}
			}(url)
		}

		var sortedResults []result
		for range urls {
			sortedResults = append(sortedResults, <-results)
		}

		sort.Slice(sortedResults, func(i, j int) bool {
			return sortedResults[i].size < sortedResults[j].size
		})

		for _, r := range sortedResults {
			if r.error != nil {
				fmt.Printf("%s: error: %s\n", r.url, r.error)
			} else {
				fmt.Printf("%s: %d bytes\n", r.url, r.size)
			}
		}

		return nil
	}

	app.Run(strings.Split("url-size url1 url2 url3 ...", " "))
}
