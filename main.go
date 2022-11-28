package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	url := flag.String("u", "https://duckduckgo.com", "enter the site url")
	flag.Parse()
	response := siteresponse(*url)

	if response != 200 {
		fmt.Println("something went wrong.")
		// do stuff
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
