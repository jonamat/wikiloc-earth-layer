package init

import (
	"encoding/xml"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
	vp "github.com/spf13/viper"
	"github.com/twpayne/go-kml/v2"
)

var (
	stopDelay       = vp.GetFloat64("stopDelay")
	serverURL       = vp.GetString("serverURL")
	updatesEndpoint = vp.GetString("endpoints.updates")
	legendEndpoint  = vp.GetString("endpoints.legend")
)

func Handle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var legendURL = fmt.Sprintf("%s%s?text=%s", serverURL, legendEndpoint, url.QueryEscape("Loading..."))
	var updatesURL = serverURL + updatesEndpoint

	k := kml.KML(
		kml.Folder(
			kml.Name("Wikiloc"),

			// Legend
			&kml.CompoundElement{
				StartElement: xml.StartElement{Name: xml.Name{Local: "ScreenOverlay"}, Attr: []xml.Attr{{Name: xml.Name{Local: "id"}, Value: "legend"}}},
				Children: []kml.Element{
					kml.Name("Trail count"),
					kml.Visibility(true),
					kml.Color(color.RGBA{255, 255, 255, 255}),
					kml.Icon(
						kml.Href(legendURL),
					),
					kml.OverlayXY(kml.Vec2{X: 0, Y: 0, XUnits: kml.UnitsFraction, YUnits: kml.UnitsFraction}),
					kml.ScreenXY(kml.Vec2{X: 10, Y: 25, XUnits: kml.UnitsPixels, YUnits: kml.UnitsPixels}),
				},
			},

			// Root folder for trails
			&kml.CompoundElement{
				StartElement: xml.StartElement{Name: xml.Name{Local: "Folder"}, Attr: []xml.Attr{{Name: xml.Name{Local: "id"}, Value: "trails"}}},
				Children: []kml.Element{
					kml.Name("Trails"),
					kml.Visibility(true),
					kml.Open(false),

					// Style definitions for the trails
					kml.SharedStyle(
						"trail",
						kml.LabelStyle(
							kml.Scale(0),
						),
						kml.LineStyle(
							kml.Color(color.RGBA{255, 255, 255, 180}),
							kml.ColorMode(kml.ColorModeRandom),
							kml.Width(2.5),
						),
					),
				},
			},

			kml.NetworkLink(
				kml.Name("Network updates"),

				kml.Visibility(true),
				kml.Open(false),
				kml.RefreshVisibility(true),
				kml.FlyToView(false),

				kml.Link(
					kml.Href(updatesURL),
					kml.ViewRefreshMode(kml.ViewRefreshModeOnStop),
					kml.ViewRefreshTime(stopDelay),
					kml.ViewFormat("view=[bboxWest],[bboxSouth],[bboxEast],[bboxNorth]"),
				),
			),
		),
	)
	w.Header().Set("content-type", "application/vnd.google-earth.kml+xml")
	if err := k.WriteIndent(w, "", "  "); err != nil {
		log.Println(err)
	}
}
