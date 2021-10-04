package updates

type Thumbnail struct {
	ID  uint64 `json:"id"`
	URL string `json:"url"`
}

type Trail struct {
	ID              uint64      `json:"id"`
	TrailRank       uint8       `json:"trailrank"`
	Distance        string      `json:"distance"`
	Elevation       string      `json:"slope"`
	Name            string      `json:"name"`
	Near            string      `json:"near"`
	Author          string      `json:"author"`
	AuthorAvatar    string      `json:"authorAvatar"`
	TrailTypeText   string      `json:"pictoText"`
	TrailType       uint8       `json:"trailType"`
	TrailTypeImgNum uint8       `json:"picto"`
	PrettyURL       string      `json:"prettyURL"`
	Lon             float64     `json:"lon"`
	Lat             float64     `json:"lat"`
	Thumbnails      []Thumbnail `json:"thumbs"`
}

type WikilocResponse struct {
	Count        uint64  `json:"count"`
	RoundedCount uint64  `json:"roundedCount"`
	Locale       string  `json:"locale"`
	Trails       []Trail `json:"spas"`
}

type Description struct {
	Type            string
	Rank            uint8
	Distance        string
	DistanceUnits   string
	Elevation       string
	ElevationUnits  string
	Author          string
	Link            string
	PathDescription string
	Thumbnails      []Thumbnail
}
