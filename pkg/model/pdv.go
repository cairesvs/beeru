// Package model that contains the struct for PDV represatation.
package model

import (
	"fmt"
)

// PDVSlice Slice of PDV
type PDVSlice struct {
	PDVS []*PDV `json:"pdvs,omitempty"`
}

// Point represent the geojson Point
type Point struct {
	Type        string    `json:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
}

// MultiPolygon represent the geojson Point
type MultiPolygon struct {
	Type        string          `json:"type,omitempty"`
	Coordinates [][][][]float64 `json:"coordinates,omitempty"`
}

// PDV model struct
type PDV struct {
	ID           string        `json:"id,omitempty"`
	TradingName  string        `json:"tradingName,omitempty"`
	OwnerName    string        `json:"ownerName,omitempty"`
	Document     string        `json:"document,omitempty"`
	CoverageArea *MultiPolygon `json:"coverageArea,omitempty"`
	Address      *Point        `json:"address,omitempty"`
}

func (p *PDV) String() string {
	return fmt.Sprintf("%s %s %s %s", p.ID, p.TradingName, p.OwnerName, p.Document)
}
