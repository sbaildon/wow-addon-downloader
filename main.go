package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"net/url"

	"github.com/sbaildon/wow-addon-downloader/providers"
	_ "github.com/sbaildon/wow-addon-downloader/providers/curseforge"
)

func main() {
	file, err := os.Open("./config.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url, err := url.ParseRequestURI(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		provider, err :=  providers.GetProvider(url.Hostname())
		if err != nil {
			log.Fatal(err)
		}


		fmt.Println(provider.DownloadURL(*url))
		fmt.Println(provider.GetVersion(*url))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
