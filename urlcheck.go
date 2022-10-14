package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const timeout time.Duration = 3 * time.Second

type httpResponse struct {
	url      string
	response *http.Response
	err      error
}

func asyncHTTPGets(urls []string, ch chan *httpResponse) {
	for _, url := range urls {
		go func(url string) {
			if !strings.HasPrefix(url, "http") {
				url = string("https://" + url)
			}
			resp, err := http.Get(url)
			ch <- &httpResponse{url, resp, err}
			fmt.Println("Visiting to", url)
		}(url)
	}
}

func writeLines(lines map[string]int, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	keys := make([]string, 0, len(lines))

	for key := range lines {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return lines[keys[i]] < lines[keys[j]]
	})

	for _, k := range keys {
		fmt.Fprintln(w, k, "\tResp Body Size: ", lines[k])
	}
	return w.Flush()
}

func main() {
	fmt.Println("Please locate your URL List (txt) file path ")
	var filePath string
	fmt.Scanln(&filePath)
	data, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer data.Close()

	var lines []string
	var s = make(map[string]int)
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	responseCount := 0
	ch := make(chan *httpResponse)
	go asyncHTTPGets(lines, ch)
	for responseCount != len(lines) {
		select {
		case r := <-ch:
			if err != nil {
				panic(err)
			} else {
				resp, err := http.Get(r.url)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				n := len(body)
				s[r.url] = n
			}
			responseCount++
		case <-time.After(timeout):
			os.Exit(1)
		}
	}

	if err := writeLines(s, filePath); err != nil {
		panic(err)
	}
	fmt.Println("Done! Please check the output on " + filePath)
}
