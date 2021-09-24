package networklink

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"image/color"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/jonamat/wikiloc-earth-layer/pkg/scraper"
	"github.com/julienschmidt/httprouter"
	vp "github.com/spf13/viper"
	"github.com/twpayne/go-kml/v2"
)

var (
	client              = &http.Client{Timeout: time.Duration(vp.GetInt("connectionTimeout")) * time.Second}
	distanceUnit        string
	elevationUnit       string
	legendEp            = vp.GetString("endpoints.legend")
	units               = vp.GetString("units")
	serverURL           = vp.GetString("serverURL")
	retryDelay          = time.Duration(vp.GetInt("retryDelay"))
	connAttempts        = vp.GetInt("connectionAttempts")
	descriptionTemplate = template.Must(template.ParseFiles(path.Join(vp.GetString("basepath"), "./web/templates/description.tmpl")))
)

func init() {
	switch vp.GetString("mesSys") {
	case "metric":
		distanceUnit = "Km"
		elevationUnit = "m"
	case "imperial":
		distanceUnit = "mi"
		elevationUnit = "ft"
	default:
		log.Println("Unkown measurement system")
	}
}

func sendEmtpy(err error, w http.ResponseWriter) {
	log.Println(err)
	var legendURL = fmt.Sprintf("%s%s?text=%s", serverURL, legendEp, url.QueryEscape(fmt.Sprintf("An error occurred during the request|See server logs for details")))
	w.Header().Set("content-type", "application/vnd.google-earth.kml+xml")
	kml.KML(
		kml.ScreenOverlay(
			kml.Name("Trail count overlay"),
			kml.Visibility(true),
			kml.Color(color.RGBA{255, 255, 255, 255}),
			kml.Icon(
				kml.Href(legendURL),
			),
			kml.OverlayXY(kml.Vec2{X: 0, Y: 0, XUnits: kml.UnitsFraction, YUnits: kml.UnitsFraction}),
			kml.ScreenXY(kml.Vec2{X: 10, Y: 25, XUnits: kml.UnitsPixels, YUnits: kml.UnitsPixels}),
		),
	).Write(w)
}

func Compose(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		sendEmtpy(err, w)
		return
	}

	// Get previous cookies to avoid trails re-fetch
	prevTrailsRaw := params.Get("ids")
	var prevTrails []uint64
	if len(prevTrailsRaw) < 1 {
		log.Printf("No loaded trails in the map")
	} else {
		tstr := strings.Split(prevTrailsRaw, "|")
		for _, t := range tstr {
			tint, err := strconv.ParseUint(t, 10, 64)
			if err != nil {
				log.Printf("Cannot convert trail ID %s into type uint64, error: %s", t, err.Error())
			} else {
				prevTrails = append(prevTrails, tint)
			}
		}

		log.Printf("Trails already loaded: %d\n", len(prevTrails))
	}

	// Return empty response if view param is not provided
	view := params.Get("view")
	if len(view) < 1 {
		sendEmtpy(fmt.Errorf("parameter \"view\" not provided"), w)
		return
	}

	// Get coordinates of Earth viewport
	coordinates := strings.Split(view, ",")
	if len(coordinates) < 4 {
		sendEmtpy(fmt.Errorf("incomplete coordinates list"), w)
		return
	}

	log.Printf("Received viewport coordinates:\n\nlongitude_west: %s\nlatitude_south: %s\nlongitude_east: %s\nlatitude_north: %s\n\n", coordinates[0], coordinates[1], coordinates[2], coordinates[3])

	// Create Wikiloc-like view coordinates
	sw := fmt.Sprintf("%s,%s", coordinates[1], coordinates[0])
	ne := fmt.Sprintf("%s,%s", coordinates[3], coordinates[2])

	// Compose wikiloc request URL
	getTrailsURL := fmt.Sprintf("https://www.wikiloc.com/wikiloc/find.do?event=map&to=24&sw=%s&ne=%s", sw, ne)

	log.Println("Making request to Wikiloc...")
	log.Printf("Request URL: %s", getTrailsURL)

	/* -------------------------------------------------------------------------- */
	/*                         Retrieve data from Wikiloc                         */
	/* -------------------------------------------------------------------------- */

	req, _ := http.NewRequest("GET", getTrailsURL, nil)
	req.Header.Add("referer", getTrailsURL)
	req.Header.Add("accept-language", "en;q=0.9")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")

	var res *http.Response
	for c := 0; c < connAttempts; c++ {
		res, err = client.Do(req)
		if err == nil && res.StatusCode < 300 {
			break
		} else {
			log.Printf(`Request error, attempt n.%d of %d`, c+1, connAttempts)
			time.Sleep(retryDelay * time.Second)
		}
	}
	if err != nil {
		sendEmtpy(err, w)
		return
	}
	if res.StatusCode > 299 {
		sendEmtpy(fmt.Errorf(`Server responds with status code %d`, res.StatusCode), w)
		return
	}

	// Store Wikiloc response body
	rawBody, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		sendEmtpy(err, w)
		return
	}

	log.Println("Request completed successfully")

	log.Println("Parsing body...")

	// Parse Wikiloc response json
	var body WikilocResponse
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		sendEmtpy(err, w)
		return
	}

	// Remove trails already in map
	var trails []Trail
	for _, t := range body.Trails {
		isNew := true
		for _, pt := range prevTrails {
			if t.ID == pt {
				isNew = false
				break
			}
		}

		if isNew {
			trails = append(trails, t)
		}
	}

	var trailsCount = len(trails)
	var legendURL = fmt.Sprintf("%s%s?text=%s", serverURL, legendEp, url.QueryEscape(fmt.Sprintf("Trails found in this area: %d|New trails: %d", body.Count, trailsCount)))

	log.Printf("Body parsed successfully. Found %d trails in his area. New trails: %d.\n", body.Count, trailsCount)

	/* -------------------------------------------------------------------------- */
	/*                            Compose KML Documents                           */
	/* -------------------------------------------------------------------------- */

	log.Println("Composing response KML file...")

	// Extracted trails ids
	var trailIDs []string

	// Iterate over parsed Trails slice
	var wg sync.WaitGroup
	docsChan := make(chan kml.Element, trailsCount)
	wg.Add(trailsCount)

	for i, trail := range trails {
		trailIDs = append(trailIDs, strconv.FormatUint(trail.ID, 10))

		go func(trail Trail, i int, docsChan chan kml.Element) {
			defer wg.Done()

			// Request the html page of the current trail
			log.Printf("[%d] Fetching %s\n", i, trail.PrettyURL)

			var res *http.Response
			for c := 0; c < connAttempts; c++ {
				res, err = client.Get(fmt.Sprintf("https://www.wikiloc.com%s", trail.PrettyURL))
				if err == nil && res.StatusCode < 300 {
					break
				} else {
					log.Printf(`[%d] Request error, attempt n.%d of %d`, i, c+1, connAttempts)
					time.Sleep(retryDelay * time.Second)
				}
			}
			if err != nil {
				log.Println(fmt.Sprintf("[%d] Error, trail skipped | %s", i, err.Error()))
				return
			}
			if res.StatusCode > 299 {
				log.Println(fmt.Sprintf("[%d] Wikiloc responded with status code %d, trail skipped", i, res.StatusCode))
				return
			}

			pageStrem, err := io.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				log.Println(fmt.Sprintf("[%d] %s", i, err.Error()))
				return
			}

			log.Printf("[%d] Fetched page %s\n", i, trail.PrettyURL)

			// Scrape path geometry from the received html code
			html := string(pageStrem)
			pathGeometry, err := scraper.GetGeometry(&html)
			if err != nil {
				log.Println(fmt.Sprintf("[%d] %s", i, err.Error()))
				return
			}

			// Create a slice of KML <coordinate> elements from the scraped geometry
			var kmlCoords []kml.Coordinate
			for _, t := range pathGeometry {
				kmlCoords = append(kmlCoords, kml.Coordinate{Lon: t.Lon, Lat: t.Lat, Alt: 0})
			}

			// Parse & convert units from imperial to metric
			distance, err := strconv.ParseFloat(trail.Distance, 64)
			if err != nil {
				distance = 0.0
			}
			if units == "metric" {
				distance = distance * 1.60934 // mi to km
			}

			elevation, err := strconv.ParseFloat(trail.Elevation, 64)
			if err != nil {
				distance = 0.0
			}
			if units == "metric" {
				elevation = elevation * 0.3048 // ft to m
			}

			// Compose the description viewport from template
			descrData := &Description{
				Type:           trail.TrailTypeText,
				Rank:           trail.TrailRank,
				Distance:       fmt.Sprintf("%.2f", distance),
				Elevation:      fmt.Sprintf("%.2f", elevation),
				Author:         trail.Author,
				Link:           fmt.Sprintf("https://www.wikiloc.com%s", trail.PrettyURL),
				Thumbnails:     trail.Thumbnails,
				DistanceUnits:  distanceUnit,
				ElevationUnits: elevationUnit,
			}

			var descrBuff bytes.Buffer
			if err := descriptionTemplate.Execute(&descrBuff, descrData); err != nil {
				log.Println(fmt.Sprintf("[%d] %s", i, err.Error()))
				return
			}
			// todo bufio.Scanner: token too long
			descr, err := io.ReadAll(&descrBuff)
			if err != nil {
				log.Println(fmt.Sprintf("[%d] %s", i, err.Error()))
				return
			}

			// Create the URL for the trail icon
			icon := fmt.Sprintf("%s/static/icons/%d.png", serverURL, trail.TrailTypeImgNum)

			// Create a KML <Document> element for the trail and append it to docs slice
			docsChan <- kml.Document(
				kml.Name(trail.Name),

				// Placemark for the icon (starting point)
				kml.Placemark(
					kml.Name(trail.Name),
					kml.Description(string(descr)),
					kml.StyleURL("#trail"),
					kml.Style(
						kml.IconStyle(
							kml.Scale(1.2),
							kml.Icon(
								kml.Href(icon),
							),
						),
					),

					kml.Point(
						kml.Coordinates(
							kml.Coordinate{Lon: trail.Lon, Lat: trail.Lat},
						),
					),
				),

				// Placemark for the path
				kml.Placemark(
					kml.Name(trail.Name),
					kml.Description(string(descr)),
					kml.StyleURL("#trail"),

					kml.LineString(
						kml.Tessellate(true),
						kml.AltitudeMode(kml.AltitudeModeClampToGround),
						kml.Coordinates(
							kmlCoords...,
						),
					),
				),
			)
		}(trail, i, docsChan)
	}

	wg.Wait()
	close(docsChan)

	var docs []kml.Element
	for doc := range docsChan {
		docs = append(docs, doc)
	}

	/* -------------------------------------------------------------------------- */
	/*                             Compose Updates KML                            */
	/* -------------------------------------------------------------------------- */

	// Append new trails to trials cookie
	cookie := strings.Join(trailIDs, "|") + "|" + prevTrailsRaw

	kmlRes := kml.KML(
		kml.NetworkLinkControl(
			kml.Cookie("ids="+cookie),

			kml.Update(
				// todo remove hc
				kml.TargetHref("http://localhost:3000/api/v1/init"),
				kml.Create(
					&kml.CompoundElement{
						StartElement: xml.StartElement{Name: xml.Name{Local: "Folder"}, Attr: []xml.Attr{{Name: xml.Name{Local: "targetId"}, Value: "trails"}}},
						Children:     docs,
					},
				),
				kml.Change(
					&kml.CompoundElement{
						StartElement: xml.StartElement{Name: xml.Name{Local: "ScreenOverlay"}, Attr: []xml.Attr{{Name: xml.Name{Local: "targetId"}, Value: "legend"}}},
						Children: []kml.Element{
							kml.Icon(
								kml.Href(legendURL),
							),
						},
					},
				),
			),
		),
	)

	log.Println("KML file composed successfully")

	log.Println("Sending response...")

	// Set content-type and send the composed KML
	w.Header().Set("content-type", "application/vnd.google-earth.kml+xml")
	if err := kmlRes.WriteIndent(w, "", "  "); err != nil {
		sendEmtpy(err, w)
		return
	}

	log.Println("Response sent successfully")
}
