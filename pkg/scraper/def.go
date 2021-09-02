package scraper

type MapData struct {
	SpaId        uint32  `json:"spaId"`
	Nom          string  `json:"nom"`
	PrettyURL    string  `json:"prettyURL"`
	Geom         string  `json:"geom"`
	ImgOnTooltip string  `json:"imgOnTooltip"`
	Blat         float32 `json:"blat"`
	Blng         float32 `json:"blng"`
	Elat         float32 `json:"elat"`
	Elng         float32 `json:"elng"`
	Loop         bool    `json:"loop"`
	Waypoint     bool    `json:"waypoint"`
}

type Data struct {
	MapData []MapData `json:"mapData"`
}

type Coordinates struct {
	Lat float64
	Lon float64
}
