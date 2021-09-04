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

// Coordinate is the base position: it differs from a point in that in contains no further
// information apart from dims and values
type Coordinate []float64

// Geometry is the base type for all twkb representations
type Geometry interface {
	Type() GeomType
	Dim() uint8
}

// Hdr represents the metadata information about a geometry
type Hdr struct {
	gtype GeomType
	dim   uint8
	bbox  BBox
}

func (h *Hdr) Type() GeomType {
	return h.gtype
}
func (h *Hdr) Dim() uint8 {
	return h.dim
}

type Point struct {
	*Hdr
	Coord Coordinate
}

type LineString struct {
	*Hdr
	Coords []Coordinate
}

type Polygon struct {
	*Hdr
	Rings []LinearRing
}

type LinearRing []Coordinate

type MultiPoint struct {
	*Hdr
	Ids    []int64
	Points []Point
}

type MultiLineString struct {
	*Hdr
	Ids         []int64
	LineStrings []LineString
}

type MultiPolygon struct {
	*Hdr
	Ids      []int64
	Polygons []Polygon
}

type GeometryCollection struct {
	*Hdr
	Ids        []int64
	Geometries []Geometry
}
