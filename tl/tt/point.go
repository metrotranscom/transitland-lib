package tt

import (
	"database/sql/driver"
	"io"

	geom "github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkb"
	"github.com/twpayne/go-geom/encoding/wkbcommon"
)

// Point is an EWKB/SL encoded point
type Point struct {
	Valid bool
	geom.Point
}

// NewPoint returns a Point from lon, lat
func NewPoint(lon, lat float64) Point {
	g := geom.NewPointFlat(geom.XY, geom.Coord{lon, lat})
	if g == nil {
		return Point{}
	}
	g.SetSRID(4326)
	return Point{Point: *g, Valid: true}
}

func (g Point) Value() (driver.Value, error) {
	if !g.Valid {
		return nil, nil
	}
	return wkbEncode(&g.Point)
}

func (g *Point) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	b, ok := src.([]byte)
	if !ok {
		return wkb.ErrExpectedByteSlice{Value: src}
	}
	// Parse
	var p geom.T
	var err error
	p, err = wkbDecode(b)
	if err != nil {
		return err
	}
	p1, ok := p.(*geom.Point)
	if !ok {
		return wkbcommon.ErrUnexpectedType{Got: p1, Want: g}
	}
	g.Valid = true
	g.Point = *p1
	return nil
}

func (g Point) String() string {
	a, _ := g.MarshalJSON()
	return string(a)
}

func (g Point) MarshalJSON() ([]byte, error) {
	if !g.Valid {
		return jsonNull(), nil
	}
	return geojsonEncode(&g.Point)
}

func (g *Point) UnmarshalGQL(v interface{}) error {
	return nil
}

func (g Point) MarshalGQL(w io.Writer) {
	b, _ := g.MarshalJSON()
	w.Write(b)
}
