package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

type URLs struct {
	Providers []string `json:"provider_urls"`
	Drugs     []string `json:"formulary_urls"`
}

func download(urls []string, client *http.Client) error {
	for _, url := range urls {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			return err
		}
		out, err := os.Create(path.Base(req.URL.Path))
		defer out.Close()
		if err != nil {
			return err
		}
		fmt.Printf("Downloading %s\n", url)
		if _, err := io.Copy(out, resp.Body); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	// get the url from args
	if len(os.Args) != 2 {
		fmt.Println("download-coverage-json [indexurl]")
		os.Exit(1)
	}
	indexURL := os.Args[1]
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
	err = download(urls.Providers, client)
	checkErr(err)
	err = download(urls.Drugs, client)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
