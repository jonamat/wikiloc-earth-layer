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

const requests = 80
const path = "./web/static/icons/"

func main() {
	var wg sync.WaitGroup

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	for counter := 0; counter < requests; counter++ {
		wg.Add(1)
		go func(client *http.Client, counter int, wg *sync.WaitGroup) {
			defer wg.Done()

			iconURL := fmt.Sprintf(`https://sc.wklcdn.com/wikiloc/images/pictograms/svg/%d.svg`, counter)

			log.Printf("[%d] Fetching %s\n", counter, iconURL)

			req, err := http.NewRequest("GET", iconURL, nil)
			if err != nil {
				log.Printf(`Error creating new request: %s`, err)
				return
			}
			req.Header.Add("referer", iconURL)
			req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")

			res, err := client.Do(req)
			if err != nil {
				log.Printf(`[%d] Error on request: %s`, counter, err)
				return
			}

			rawBody, err := io.ReadAll(res.Body)
			defer res.Body.Close()

			if res.StatusCode > 299 {
				log.Printf(`[%d] Server responds with status code %d for the icon %s`, counter, res.StatusCode, fmt.Sprint(counter))
				return
			}

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
			err = imgtools.SavePNG(icon, fmt.Sprint(counter), path)
			if err != nil {
				log.Printf(`[%d] Error on image creation: %s`, counter, err)
				return
			}
		}(client, counter, &wg)
	}

	wg.Wait()
}
