package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	imgtools "github.com/wikiloc-layer/pkg/img_tools"
)

const (
	requests = 80
	path     = "./web/static/icons/"
)

func main() {
	totalIcons := 0
	var wg sync.WaitGroup

	client := &http.Client{Timeout: 5 * time.Second}

	for counter := 0; counter < requests; counter++ {
		wg.Add(1)
		go func(client *http.Client, counter int, wg *sync.WaitGroup) {
			defer wg.Done()
			//wklcdn
			iconURL := fmt.Sprintf(`https://sc.wklcdn.com/wikiloc/images/pictograms/svg/%d.svg`, counter)

			log.Printf("[%d] Fetching %s\n", counter, iconURL)

			req, err := http.NewRequest("GET", iconURL, nil)
			if err != nil {
				log.Printf(`Error creating new request: %s`, err)
				return
			}
			req.Header.Add("referer", iconURL)
			req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")

			var res *http.Response
			for c := 0; c < 2; c++ {
				res, err = client.Do(req)
				if err == nil && res.StatusCode < 300 {
					break
				}
			}
			if err != nil {
				log.Printf(`[%d] Error on request: %s`, counter, err)
				return
			}
			if res.StatusCode > 299 {
				log.Printf(`[%d] Server responds with status code %d for the icon %s`, counter, res.StatusCode, fmt.Sprint(counter))
				return
			}

			rawBody, err := io.ReadAll(res.Body)
			defer res.Body.Close()

			if err != nil {
				log.Printf(`[%d] Error on body parsing: %s`, counter, err)
				return
			}

			icon, err := imgtools.MakeIcon(&rawBody)
			if err != nil {
				log.Printf(`[%d] Error on image creation: %s`, counter, err)
				return
			}

			// Save the cropped icon into file
			err = imgtools.SavePNG(icon, path, fmt.Sprint(counter))
			if err != nil {
				log.Printf(`[%d] Error on image creation: %s`, counter, err)
				return
			}

			totalIcons++
		}(client, counter, &wg)
	}

	wg.Wait()

	fmt.Printf("Successfully downloaded %d icons. Tested %d URLs\n", totalIcons, requests)
}
