package index

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "https://jonamat.github.io/wikiloc-earth-layer/", 301)
}
