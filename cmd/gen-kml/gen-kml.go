package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"

	_setup "github.com/jonamat/wikiloc-earth-layer/pkg/_setup"
	vp "github.com/spf13/viper"
	"github.com/twpayne/go-kml/v2"
)

func init() {
	// Load configuration and set viper singleton
	_setup.Init()
}

func main() {
	initEndpoint := vp.GetString("endpoints.init")
	protocol := vp.GetString("protocol")
	host := vp.GetString("host")
	port := vp.GetString("port")

	if len(protocol) == 0 || len(host) == 0 {
		panic("Protocol or host not defined")
	}

	var url string
	if port != "80" || len(port) == 0 {
		url = fmt.Sprintf("%s://%s:%s%s", protocol, host, port, initEndpoint)
	} else {
		url = fmt.Sprintf("%s://%s%s", protocol, host, initEndpoint)
	}

	kml := kml.KML(
		// First network link (will be replaced by networklinkcontrol)
		kml.NetworkLink(
			kml.Name("Wikiloc"),

			kml.Visibility(true),
			kml.Open(false),
			kml.RefreshVisibility(true),
			kml.FlyToView(false),

			kml.Link(
				kml.Href(url),
				kml.ViewRefreshMode(kml.ViewRefreshModeOnRequest),
				kml.ViewRefreshTime(0),
				kml.ViewFormat(""),
			),
		),
	)

	archive, err := os.Create("./web/static/wikiloc-earth-layer.kmz")
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)

	writer, err := zipWriter.Create("./wikiloc-earth-layer.kml")
	if err != nil {
		panic(err)
	}

	// file, _ := os.Create("./web/static/wikiloc-earth-layer.kml")
	// defer file.Close()

	if err := kml.WriteIndent(writer, "", "  "); err != nil {
		panic(err)
	}
	zipWriter.Close()

	log.Printf("Generated init KML with the following vars:\nPROTOCOL: %s\nHOST: %s\nPORT: %s", protocol, host, port)
	log.Println("Compressed and saved in ./web/static/wikiloc-earth-layer.kmz")
}
