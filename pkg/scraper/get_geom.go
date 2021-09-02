package scraper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/devork/twkb"
)

func GetGeometry(html *string) ([]Coordinates, error) {
	// Split html code into lines
	lines, err := StringToLines(html)
	if err != nil {
		return nil, err
	}

	// Search the line that contains "mapData"
	var rowData string
	for _, line := range lines {
		if strings.Contains(line, "mapData") {
			// Remove = and ;
			index := strings.Index(line, "=")
			rowData = line[index+1 : len(line)-1]
		}
	}
	if len(rowData) == 0 {
		return nil, fmt.Errorf("trail geometry not found")
	}

	// Parse json into Data struct
	var data Data
	err = json.Unmarshal([]byte(rowData), &data)
	if err != nil {
		return nil, err
	}

	// Decode parsed data from base64
	twkbGeom, err := base64.StdEncoding.DecodeString(data.MapData[0].Geom)
	if err != nil {
		return nil, err
	}

	// Decode geometry from TWKB encoding
	// See https://github.com/TWKB/Specification/blob/master/twkb.md
	geom, err := twkb.Decode(bytes.NewReader([]byte(twkbGeom)))
	if err != nil {
		return nil, err
	}

	lineString, ok := geom.(*twkb.LineString)
	if !ok {
		return nil, fmt.Errorf("trail type is not a line string")
	}

	var trail []Coordinates
	for _, point := range lineString.Coords {
		trail = append(trail, Coordinates{Lat: point[1], Lon: point[0]})
	}

	return trail, nil
}
