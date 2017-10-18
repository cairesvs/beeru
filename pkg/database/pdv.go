package database

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	"github.com/cairesvs/beeru/pkg/model"
)

type PDVGetInput struct {
	Database *BeeruDatabase
	ID       string
}

func GetPDV(input *PDVGetInput) *model.PDV {
	db := input.Database.DB
	rows, err := db.Query("SELECT id,trading_name, owner_name, document, ST_AsGeoJSON(coverage_area) as coverage_area, ST_AsGeoJSON(address) as address FROM pdv where id = $1", input.ID)
	if err != nil {
		log.Fatal(err)
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
			log.Println(err)
		}
		mp := &model.MultiPolygon{}
		p := &model.Point{}
		err = json.Unmarshal([]byte(coverageArea), mp)
		if err != nil {
			log.Println("Failed to unmarshall the multipolygon geojson", err)
		}
		err = json.Unmarshal([]byte(address), p)
		if err != nil {
			log.Println("Failed to unmarshall the point geojson", err)
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

type PDVCreateInput struct {
	Database *BeeruDatabase
	PDV      *model.PDV
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func CreatePDV(input *PDVCreateInput) error {
	db := input.Database.DB
	pdv := input.PDV
	if len(pdv.Document) != 17 {
		return fmt.Errorf("The Document must have 17 digits")
	}
	bytesCoverageArea, _ := json.Marshal(pdv.CoverageArea)
	bytesAddress, _ := json.Marshal(pdv.Address)
	query := fmt.Sprintf("INSERT INTO pdv(trading_name,owner_name,document,coverage_area, address) VALUES($1,$2,$3, ST_SetSRID(ST_GeomFromGeoJSON('%s'), 4326),ST_SetSRID(ST_GeomFromGeoJSON('%s'), 4326))", string(bytesCoverageArea), string(bytesAddress))
	_, err := db.Exec(query, pdv.TradingName, pdv.OwnerName, RandStringBytes(17))
	if err != nil {
		return fmt.Errorf("The query is malformed %s", err)
	}
	return nil
}

type PDVFindInput struct {
	Database *BeeruDatabase
	Point    *model.Point
}

func FindPDV(input *PDVFindInput) *model.PDVSlice {
	db := input.Database.DB
	rows, err := db.Query(fmt.Sprintf("SELECT id,trading_name, owner_name, document, ST_AsGeoJSON(coverage_area) as coverage_area, ST_AsGeoJSON(address) as address FROM pdv WHERE ST_Contains(pdv.coverage_area, ST_GeomFromText('POINT(%03.6f %03.6f)', 4326)) ORDER BY ST_Distance(pdv.coverage_area, pdv.address)", input.Point.Coordinates[0], input.Point.Coordinates[1]))
	if err != nil {
		log.Fatal(err)
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
			log.Println(err)
		}
		mp := &model.MultiPolygon{}
		p := &model.Point{}
		err = json.Unmarshal([]byte(coverageArea), mp)
		if err != nil {
			log.Println("Failed to unmarshall the multipolygon geojson", err)
		}
		err = json.Unmarshal([]byte(address), p)
		if err != nil {
			log.Println("Failed to unmarshall the point geojson", err)
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
