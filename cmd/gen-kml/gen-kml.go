package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/twpayne/go-kml"
)

const overlayEndpoint = "/api/v1/overlay"

func main() {
	godotenv.Load()

	protocol := os.Getenv("PROTOCOL")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if len(protocol) == 0 || len(host) == 0 {
		log.Fatal("Protocol or host not defined")
	}

	var url string
	if port != "80" || len(port) == 0 {
		url = fmt.Sprintf("%s://%s:%s%s", protocol, host, port, overlayEndpoint)
	} else {
		url = fmt.Sprintf("%s://%s%s", protocol, host, overlayEndpoint)
	}

	kml := kml.KML(
		kml.NetworkLink(
			kml.Name("Wikiloc"),

			kml.Visibility(true),
			kml.Open(true),
			kml.RefreshVisibility(true),
			kml.FlyToView(false),

			kml.Link(
				kml.Href(url),
				kml.ViewRefreshMode(kml.ViewRefreshModeOnRequest),
				kml.ViewRefreshTime(0),
				kml.ViewFormat("view=[bboxWest],[bboxSouth],[bboxEast],[bboxNorth]"),
			),
		),
	)

	file, err := os.Create("./web/static/wikiloc-earth-layer.kml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := kml.WriteIndent(file, "", "  "); err != nil {
		panic(err)
	}
}
