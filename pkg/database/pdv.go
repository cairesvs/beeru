// Package database repository for PDV.
//
// Has the 3 basic functionalities: Get, Create and Find.
package database

import (
	"encoding/json"
	"fmt"

	"github.com/cairesvs/beeru/pkg/logger"
	"github.com/cairesvs/beeru/pkg/model"
)

// PDVGetInput Input to get PDV for given ID
type PDVGetInput struct {
	Database *BeeruDatabase
	ID       string
}

// GetPDV get PDV for given ID
func GetPDV(input *PDVGetInput) *model.PDV {
	db := input.Database.DB
	rows, err := db.Query("SELECT id,trading_name, owner_name, document, ST_AsGeoJSON(coverage_area) as coverage_area, ST_AsGeoJSON(address) as address FROM pdv where id = $1", input.ID)
	if err != nil {
		logger.Fatal(err)
	}
	var pdv *model.PDV
	for rows.Next() {
		var id string
		var trandingName string
		var ownerName string
		var document string
		var coverageArea string
		var address string
		err = rows.Scan(&id, &trandingName, &ownerName, &document, &coverageArea, &address)
		if err != nil {
			logger.Error(err)
		}
		mp := &model.MultiPolygon{}
		p := &model.Point{}
		err = json.Unmarshal([]byte(coverageArea), mp)
		if err != nil {
			logger.Errorf("Failed to unmarshall the multipolygon geojson %s", err)
		}
		err = json.Unmarshal([]byte(address), p)
		if err != nil {
			logger.Errorf("Failed to unmarshall the point geojson %s", err)
		}
		pdv = &model.PDV{
			ID:           id,
			TradingName:  trandingName,
			OwnerName:    ownerName,
			Document:     document,
			CoverageArea: mp,
			Address:      p,
		}
	}
	return pdv
}

// PDVCreateInput input for create new PDV
type PDVCreateInput struct {
	Database *BeeruDatabase
	PDV      *model.PDV
}

// CreatePDV creates PDV on database
func CreatePDV(input *PDVCreateInput) (int64, error) {
	db := input.Database.DB
	pdv := input.PDV
	if len(pdv.Document) != 17 {
		return 0, fmt.Errorf("The Document must have 17 digits")
	}
	bytesCoverageArea, _ := json.Marshal(pdv.CoverageArea)
	bytesAddress, _ := json.Marshal(pdv.Address)
	query := fmt.Sprintf("INSERT INTO pdv(trading_name,owner_name,document,coverage_area, address) VALUES($1,$2,$3, ST_SetSRID(ST_GeomFromGeoJSON('%s'), 4326),ST_SetSRID(ST_GeomFromGeoJSON('%s'), 4326)) RETURNING id", string(bytesCoverageArea), string(bytesAddress))
	var id int64
	err := db.QueryRow(query, pdv.TradingName, pdv.OwnerName, pdv.Document).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("The query is malformed %s", err)
	}
	return id, nil
}

// PDVFindInput Input for find PDVs for given point
type PDVFindInput struct {
	Database *BeeruDatabase
	Point    *model.Point
}

// FindPDV Find PDVs for given point
func FindPDV(input *PDVFindInput) *model.PDVSlice {
	db := input.Database.DB
	rows, err := db.Query(fmt.Sprintf("SELECT id,trading_name, owner_name, document, ST_AsGeoJSON(coverage_area) as coverage_area, ST_AsGeoJSON(address) as address FROM pdv WHERE ST_Contains(pdv.coverage_area, ST_GeomFromText('POINT(%03.6f %03.6f)', 4326)) ORDER BY ST_Distance(pdv.coverage_area, pdv.address)", input.Point.Coordinates[0], input.Point.Coordinates[1]))
	if err != nil {
		logger.Error(err)
		return &model.PDVSlice{}
	}
	pdvs := []*model.PDV{}
	for rows.Next() {
		var id string
		var trandingName string
		var ownerName string
		var document string
		var coverageArea string
		var address string
		err = rows.Scan(&id, &trandingName, &ownerName, &document, &coverageArea, &address)
		if err != nil {
			logger.Error(err)
			return &model.PDVSlice{}
		}
		mp := &model.MultiPolygon{}
		p := &model.Point{}
		err = json.Unmarshal([]byte(coverageArea), mp)
		if err != nil {
			logger.Errorf("Failed to unmarshall the multipolygon geojson %s", err)
			return &model.PDVSlice{}
		}
		err = json.Unmarshal([]byte(address), p)
		if err != nil {
			logger.Errorf("Failed to unmarshall the point geojson %s", err)
			return &model.PDVSlice{}
		}
		pdv := &model.PDV{
			ID:           id,
			TradingName:  trandingName,
			OwnerName:    ownerName,
			Document:     document,
			CoverageArea: mp,
			Address:      p,
		}
		pdvs = append(pdvs, pdv)
	}
	return &model.PDVSlice{
		PDVS: pdvs,
	}

}
