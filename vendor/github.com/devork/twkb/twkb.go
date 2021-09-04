/*
Copyright [2015] Alex Davies-Moore

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package twkb

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

// BBox is represents the bounding box for the geom. The format is <minx miny (minz|minm minz) maxx maxy (maxz|maxm maxz)>
type BBox []float64

// GeomType represents the specific geometry type
type GeomType uint8

// Standard geometry types
const (
	_ GeomType = iota
	POINT
	LINESTRING
	POLYGON
	MULTIPOINT
	MULTILINESTRING
	MULTIPOLYGON
	GEOMETRYCOLLECTION
)

// Dimension types
const (
	XY uint8 = iota
	XYZ
	XYM
	XYZM
)

// masks
const (
	bbox   uint8 = 0x01
	size   uint8 = 0x02
	idlist uint8 = 0x04
	exprec uint8 = 0x08
	empty  uint8 = 0x10
)

// Common errors
var (
	ErrUnknownGeom = errors.New("unknown geometry type")
)

type decoder struct {
	r                         io.Reader
	refpoint                  [4]int64
	factors                   [4]int64
	ndims                     int
	bbox, size, idlist, empty bool
	length                    uint64
}

type mhdr struct {
	idlist []int64
	ngeoms int
}

func (d *decoder) read(v interface{}) error {
	return binary.Read(d.r, binary.BigEndian, v)
}

// Decode converts the given data into a Geometry type
func Decode(r io.Reader) (Geometry, error) {
	var flag uint8
	var decoder = &decoder{r: r}

	err := decoder.read(&flag)

	if err != nil {
		return nil, err
	}

	hdr := Hdr{}

	// Geom type and precision header
	hdr.gtype = GeomType(flag & 0x0F)

	if hdr.gtype > GEOMETRYCOLLECTION || hdr.gtype < POINT {
		return nil, ErrUnknownGeom
	}

	precision := unzigzag(uint64(flag&0xF0)) >> 4
	decoder.factors[0] = int64(math.Pow10(int(precision)))
	decoder.factors[1] = decoder.factors[0]

	// Metadata
	err = decoder.read(&flag)

	if err != nil {
		return nil, err
	}

	decoder.bbox = (flag&bbox == bbox)
	decoder.size = (flag&size == size)
	decoder.idlist = (flag&idlist == idlist)
	decoder.empty = (flag&empty == empty)

	hdr.dim = XY
	decoder.ndims = 2

	// Z and/or M Dims
	if flag&exprec == exprec {
		var extended uint8

		err := decoder.read(&extended)

		if err != nil {
			return nil, err
		}

		switch {
		case extended&XYZM == XYZM:
			hdr.dim = XYZM
			decoder.ndims = 4
			zprec := int((extended & 0x1C) >> 2)
			mprec := int((extended & 0xE0) >> 5)
			decoder.factors[2] = int64(math.Pow10(zprec))
			decoder.factors[3] = int64(math.Pow10(mprec))

		case extended&XYZ == XYZ:
			hdr.dim = XYZ
			decoder.ndims = 3
			zprec := int((extended & 0x1C) >> 2)
			decoder.factors[2] = int64(math.Pow10(zprec))

		case extended&XYM == XYM:
			hdr.dim = XYM
			decoder.ndims = 3
			mprec := int((extended & 0xE0) >> 5)
			decoder.factors[2] = int64(math.Pow10(mprec))

		}
	}

	if decoder.size {
		decoder.length, err = readVarInt64(decoder)

		if err != nil {
			return nil, err
		}
	}

	if decoder.bbox {
		var bbox = make(BBox, decoder.ndims*2)
		var min, max int64

		for idx := 0; idx < decoder.ndims; idx++ {
			min, err = readVarSInt64(decoder)

			if err != nil {
				return nil, err
			}

			max, err = readVarSInt64(decoder)

			if err != nil {
				return nil, err
			}

			bbox[idx] = float64(min)
			bbox[idx+decoder.ndims] = float64(min + max)
		}

		hdr.bbox = bbox

	}

	switch hdr.gtype {
	case POINT:
		return decodePoint(decoder, &hdr)
	case LINESTRING:
		return decodeLineString(decoder, &hdr)
	case POLYGON:
		return decodePolygon(decoder, &hdr)
	case MULTIPOINT:
		return decodeMultiPoint(decoder, &hdr)
	case MULTILINESTRING:
		return decodeMultiLineString(decoder, &hdr)
	case MULTIPOLYGON:
		return decodeMultiPolygon(decoder, &hdr)
	case GEOMETRYCOLLECTION:
		return decodeCollection(decoder, &hdr)
	}

	return &hdr, nil
}

func decodePoint(d *decoder, h *Hdr) (*Point, error) {
	coords, bbox, err := readCoords(d, 1)

	if err != nil {
		return nil, err
	}

	if h != nil {

		if bbox != nil {
			h.bbox = bbox
		}
		return &Point{h, coords[0]}, nil
	}

	return &Point{nil, coords[0]}, nil
}

func decodeLineString(d *decoder, h *Hdr) (*LineString, error) {
	size, err := readVarInt64(d)

	if err != nil {
		return nil, err
	}

	coords, bbox, err := readCoords(d, int(size))

	if err != nil {
		return nil, err
	}

	if h != nil {
		if bbox != nil {
			h.bbox = bbox
		}

		return &LineString{h, coords}, nil
	}

	return &LineString{nil, coords}, nil
}

func decodePolygon(d *decoder, h *Hdr) (*Polygon, error) {
	size, err := readVarInt64(d)

	if err != nil {
		return nil, err
	}

	var ringSize uint64
	rings := make([]LinearRing, size)
	for ring := 0; ring < int(size); ring++ {
		ringSize, err = readVarInt64(d)

		if err != nil {
			return nil, err
		}

		rings[ring], _, err = readCoords(d, int(ringSize))
	}

	if h != nil {
		return &Polygon{h, rings}, nil
	}

	return &Polygon{nil, rings}, nil
}

func decodeMultiPoint(d *decoder, h *Hdr) (*MultiPoint, error) {
	var err error
	m, err := readMHdr(d)

	if err != nil {
		return nil, err
	}

	var points = make([]Point, m.ngeoms)
	var p *Point
	for idx := 0; idx < m.ngeoms; idx++ {
		p, err = decodePoint(d, nil)

		if err != nil {
			return nil, err
		}

		points[idx] = *p
	}

	return &MultiPoint{h, m.idlist, points}, nil
}

func decodeMultiLineString(d *decoder, h *Hdr) (*MultiLineString, error) {
	var err error
	m, err := readMHdr(d)

	if err != nil {
		return nil, err
	}

	var lines = make([]LineString, m.ngeoms)
	var p *LineString
	for idx := 0; idx < m.ngeoms; idx++ {
		p, err = decodeLineString(d, nil)

		if err != nil {
			return nil, err
		}

		lines[idx] = *p
	}

	return &MultiLineString{h, m.idlist, lines}, nil
}

func decodeMultiPolygon(d *decoder, h *Hdr) (*MultiPolygon, error) {
	var err error
	m, err := readMHdr(d)

	if err != nil {
		return nil, err
	}

	var polys = make([]Polygon, m.ngeoms)
	var p *Polygon
	for idx := 0; idx < m.ngeoms; idx++ {
		p, err = decodePolygon(d, nil)

		if err != nil {
			return nil, err
		}

		polys[idx] = *p
	}

	return &MultiPolygon{h, m.idlist, polys}, nil
}

func decodeCollection(d *decoder, h *Hdr) (*GeometryCollection, error) {
	var err error
	m, err := readMHdr(d)

	if err != nil {
		return nil, err
	}

	var geoms = make([]Geometry, m.ngeoms)
	var g Geometry

	for idx := 0; idx < m.ngeoms; idx++ {
		g, err = Decode(d.r)

		if err != nil {
			return nil, err
		}

		geoms[idx] = g
	}

	return &GeometryCollection{h, m.idlist, geoms}, nil
}

func readCoords(d *decoder, size int) ([]Coordinate, BBox, error) {
	var coords = make([]Coordinate, size)
	var coord Coordinate
	var val int64
	var err error
	var bbox BBox

	if !d.bbox {
		bbox = make(BBox, d.ndims*2)
		for idx := 0; idx < d.ndims; idx++ {
			bbox[idx] = math.MaxInt64
			bbox[idx+d.ndims] = math.MinInt64
		}
	}
	for i := 0; i < size; i++ {
		coord = make(Coordinate, d.ndims)
		for j := 0; j < d.ndims; j++ {
			val, err = readVarSInt64(d)

			if err != nil {
				return nil, nil, err
			}

			d.refpoint[j] += val
			coord[j] = float64(d.refpoint[j]) / float64(d.factors[j])

			if bbox != nil {
				if coord[j] < bbox[j] {
					bbox[j] = coord[j]
				}

				if coord[j] > bbox[j+d.ndims] {
					bbox[j+d.ndims] = coord[j]
				}
			}
		}
		coords[i] = coord

	}

	return coords, bbox, nil
}

func unzigzag(nVal uint64) int64 {
	if (nVal & 1) == 0 {
		return int64(nVal >> 1)
	}
	return int64(-(nVal >> 1) - 1)
}

func readVarInt64(d *decoder) (uint64, error) {
	var byte uint8
	var val uint64
	var shift uint
	var err error
	for {
		err = d.read(&byte)

		if err != nil {
			return 0, err
		}

		if byte&0x80 == 0 {
			return val | (uint64(byte) << shift), nil
		}

		val = val | (uint64(byte)&0x7f)<<shift
		shift += 7
	}
}

func readVarSInt64(d *decoder) (int64, error) {
	val, err := readVarInt64(d)

	if err != nil {
		return 0, err
	}

	return unzigzag(val), nil
}

func readIDList(d *decoder, ngeoms int64) ([]int64, error) {
	var idlist = make([]int64, ngeoms)
	var id int64
	var err error

	for idx := 0; idx < int(ngeoms); idx++ {
		id, err = readVarSInt64(d)

		if err != nil {
			return nil, err
		}

		idlist[idx] = id
	}

	return idlist, nil
}

func readMHdr(d *decoder) (*mhdr, error) {
	var m = &mhdr{}

	ngeoms, err := readVarInt64(d)

	if err != nil {
		return nil, err
	}

	m.ngeoms = int(ngeoms)

	if d.idlist {
		m.idlist, err = readIDList(d, int64(ngeoms))

		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

//https://github.com/golang/go/issues/4594

func round(f float64) int {
	if math.Abs(f) < 0.5 {
		return 0
	}
	return int(f + math.Copysign(0.5, f))
}
