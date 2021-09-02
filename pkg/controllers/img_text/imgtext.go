package imgtext

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

func GenerateImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Println(err)
		http.Error(w, "Cannot handle request", 400)
		return
	}

	text := params.Get("text")
	if len(text) < 1 {
		log.Println("text param not provided")
		http.Error(w, "Cannot handle request", 400)
		return
	}

	log.Printf("Converting text \"%s\" into image...", text)

	res, err := http.Get("https://chart.apis.google.com/chart?chst=d_text_outline&chld=000000|20|l|ffffff|b|" + url.QueryEscape(text))

	if err != nil {
		log.Println(err)
		http.Error(w, "Cannot handle request", 400)
		return
	}

	image, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		log.Printf("img-gen: server responded with error code: %d", res.StatusCode)
		http.Error(w, "Internal Error", 500)
		return
	}
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Error", 500)
		return
	}

	w.Write(image)

	log.Println("Image sent succefully")
}
