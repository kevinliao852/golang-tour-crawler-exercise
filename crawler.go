package main

import (
	"fmt"
	"sync"
	//"time"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// 	Data := make(map[string]boolean)

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {

	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	ch := make(chan string)
	visited := make(map[string]bool)
	quit := make(chan int)
	var wg sync.WaitGroup
	var mu sync.Mutex

	go func() {
		var crawl func(url string, depth int, fetcher Fetcher)
		wg.Add(1)
		crawl = func(url string, depth int, fetcher Fetcher) {
			defer wg.Done()
			fmt.Println(wg)

			if depth <= 0 {
				return
			}

			if visited[url] == false {
				mu.Lock()
				visited[url] = true
				mu.Unlock()

				_, urls, err := fetcher.Fetch(url)
				if err != nil {
					fmt.Println(err)
					return
				}

				for _, u := range urls {
					wg.Add(1)
					go crawl(u, depth-1, fetcher)
				}
				ch <- url
			}
		}

		crawl(url, depth, fetcher)

		wg.Wait()

		quit <- 1

	}()

	for {
		select {
		case x := <-ch:
			fmt.Println(x)
		case <-quit:
			return
		}
	}

	return
}

func main() {
	Crawl("https://golang.org/", 3, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

