package index

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
)

var redirectPage = viper.GetString("redirectPage")

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, redirectPage, http.StatusMovedPermanently)
}
