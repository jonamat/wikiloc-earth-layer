package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	_setup "github.com/jonamat/wikiloc-earth-layer/pkg/_setup"
	imgtext "github.com/jonamat/wikiloc-earth-layer/pkg/controllers/imgtext"
	"github.com/jonamat/wikiloc-earth-layer/pkg/controllers/index"
	init_ctr "github.com/jonamat/wikiloc-earth-layer/pkg/controllers/init"
	updates "github.com/jonamat/wikiloc-earth-layer/pkg/controllers/updates"
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
	_setup.Init()
}

func main() {
	// Define routes
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("./web/static"))
	router.GET("/", index.Handle)
	router.GET(vp.GetString("endpoints.init"), init_ctr.Handle)
	router.GET(vp.GetString("endpoints.updates"), updates.Handle)
	router.GET(vp.GetString("endpoints.legend"), imgtext.Handle)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", vp.GetString("servicePort")))
	if err != nil {
		panic(err)
	}

	log.Printf("Server started on port %s", vp.GetString("servicePort"))

	if err := http.Serve(listener, logger(router)); err != nil {
		panic(err)
	}

}
