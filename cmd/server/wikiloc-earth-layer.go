package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	imgtext "github.com/wikiloc-layer/pkg/controllers/img_text"
	networklink "github.com/wikiloc-layer/pkg/controllers/network_link"
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

const OverlayEp = "/api/v1/overlay"
const GenImgEp = "/api/v1/gen-img"

func main() {
	godotenv.Load()
	protocol := os.Getenv("PROTOCOL")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	// Check envs
	if len(protocol) == 0 || len(host) == 0 {
		log.Fatal("Undefined protocol or host")
	}

	// Define server URL
	var url string
	if port != "80" || len(port) == 0 {
		url = fmt.Sprintf("%s://%s:%s", protocol, host, port)
	} else {
		url = fmt.Sprintf("%s://%s", protocol, host)
	}
	os.Setenv("URL", url)

	// Define paths
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("./web/static"))
	router.GET(OverlayEp, networklink.Compose)
	router.GET(GenImgEp, imgtext.GenerateImage)

	// Use 80 as default port
	if len(port) == 0 {
		port = "80"
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		panic(err)
	}

	log.Printf("Server started on port %s", port)

	if err := http.Serve(listener, logger(router)); err != nil {
		panic(err)
	}

}
