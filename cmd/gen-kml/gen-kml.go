package main

import (
	"fmt"
	"log"
	"os"

	vp "github.com/spf13/viper"
	"github.com/twpayne/go-kml"
	_init "github.com/wikiloc-layer/pkg/_init"
)

func init() {
	// Load configuration and set viper singleton
	_init.Init()
}

func main() {
	overlayEndpoint := vp.GetString("endpoints.updates")
	protocol := vp.GetString("protocol")
	host := vp.GetString("host")
	port := vp.GetString("port")

	if len(protocol) == 0 || len(host) == 0 {
		panic("Protocol or host not defined")
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

	log.Printf("Generated init KML with the following vars:\nPROTOCOL: %s\nHOST: %s\nPORT: %s", protocol, host, port)
}
