package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cretz/bine/tor"
	"golang.org/x/net/html"
)

func main() {
	url := flag.String("u", "https://duckduckgo.com", "enter the site url")
	flag.Parse()
	response := siteresponse(*url)

	if response != 200 {
		fmt.Println("something went wrong.")
	} else {
		fmt.Printf("Starting tor and fetching title of %s, please wait a few seconds...", *url)
		t, err := tor.Start(nil, nil)
		if err != nil {
			fmt.Println(err)
		}
		defer t.Close()
		// Wait at most a minute to start network and get
		dialCtx, dialCancel := context.WithTimeout(context.Background(), time.Minute)
		defer dialCancel()
		// Make connection
		dialer, err := t.Dialer(dialCtx, nil)
		if err != nil {
			fmt.Println(err)
		}
		httpClient := &http.Client{Transport: &http.Transport{DialContext: dialer.DialContext}}
		// Get /
		resp, err := httpClient.Get(*url)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		// Grab the <title>
		parsed, err := html.Parse(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Title: %v\n", getTitle(parsed))
	}

	fmt.Println("site url: ", *url)
}

func siteresponse(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	return resp.StatusCode
}

func getTitle(n *html.Node) string {
	if n.Type == html.ElementNode && n.Data == "title" {
		var title bytes.Buffer
		if err := html.Render(&title, n.FirstChild); err != nil {
			panic(err)
		}
		return strings.TrimSpace(title.String())
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if title := getTitle(c); title != "" {
			return title
		}
	}
	return ""
}
