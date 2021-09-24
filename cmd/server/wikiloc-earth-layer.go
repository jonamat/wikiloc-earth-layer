package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	_init "github.com/jonamat/wikiloc-earth-layer/pkg/_init"
	imgtext "github.com/jonamat/wikiloc-earth-layer/pkg/controllers/img_text"
	"github.com/jonamat/wikiloc-earth-layer/pkg/controllers/index"
	init_updates "github.com/jonamat/wikiloc-earth-layer/pkg/controllers/init"
	networklink "github.com/jonamat/wikiloc-earth-layer/pkg/controllers/updates"
	"github.com/julienschmidt/httprouter"
	vp "github.com/spf13/viper"
)

type Middleware struct {
	next http.Handler
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming request. Method: %s. URI: %s", r.Method, r.RequestURI)
	m.next.ServeHTTP(w, r)
}

func logger(next http.Handler) *Middleware {
	return &Middleware{next: next}
}

func init() {
	// Load configuration and set viper singleton
	_init.Init()
}

func main() {
	// Define routes
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("./web/static"))
	router.GET("/", index.Index)
	router.GET(vp.GetString("endpoints.init"), init_updates.Init_updates)
	router.GET(vp.GetString("endpoints.updates"), networklink.Compose)
	router.GET(vp.GetString("endpoints.legend"), imgtext.GenerateImage)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", vp.GetString("servicePort")))
	if err != nil {
		panic(err)
	}

	log.Printf("Server started on port %s", vp.GetString("servicePort"))

	if err := http.Serve(listener, logger(router)); err != nil {
		panic(err)
	}

}
