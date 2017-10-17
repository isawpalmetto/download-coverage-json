package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
)

type URLs struct {
	Providers []string `json:"provider_urls"`
	Drugs     []string `json:"formulary_urls"`
}

func download(client *http.Client, wg *sync.WaitGroup, urls []string, dest string) error {
	wg.Add(len(urls))
	for _, url := range urls {
		go func(u string) {
			defer wg.Done()
			req, err := http.NewRequest("GET", u, nil)
			checkErr(err)
			resp, err := client.Do(req)
			checkErr(err)
			defer resp.Body.Close()
			err = os.MkdirAll(dest, os.ModePerm)
			checkErr(err)
			out, err := os.Create(dest + path.Base(req.URL.Path))
			checkErr(err)
			defer out.Close()
			fmt.Printf("Downloading %s to %s\n", u, dest)
			_, err = io.Copy(out, resp.Body)
			checkErr(err)
		}(url)
	}
	return nil
}

func main() {
	// get the url from args
	if len(os.Args) != 3 {
		fmt.Println("download-coverage-json [indexurl] [destination]")
		os.Exit(1)
	}
	indexURL := os.Args[1]
	dest := os.Args[2]
	resp, err := http.Get(indexURL)
	defer resp.Body.Close()
	checkErr(err)
	body, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	var urls URLs
	if err := json.Unmarshal(body, &urls); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := &http.Client{}
	var wg sync.WaitGroup
	var other sync.WaitGroup
	err = download(client, &wg, urls.Providers, dest+"/providers/")
	wg.Wait()
	checkErr(err)
	err = download(client, &other, urls.Drugs, dest+"/drugs/")
	checkErr(err)
	other.Wait()

}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
