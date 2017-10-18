package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"

	"github.com/cairesvs/beeru/pkg/data"
	"github.com/cairesvs/beeru/pkg/model"
)

// BeeruDatabase represent wrapper for sql.DB
type BeeruDatabase struct {
	DB *sql.DB
}

var (
	// PDVS represent the singleton for PDVSlice
	PDVS *model.PDVSlice
	// Database represent the singleton for BeeruDatabae
	Database     *BeeruDatabase
	fileOnce     sync.Once
	databaseOnce sync.Once
)

// GetPDVS from static file
func getPDVS() *model.PDVSlice {
	fileOnce.Do(func() {
		data, err := data.Asset("data/pdvs.json")
		if err != nil {
			log.Fatal("Couldn't load base json")
		}
		pdvs := &model.PDVSlice{}
		err = json.Unmarshal(data, pdvs)
		if err != nil {
			log.Fatal("Couldn't unmarshall ", err)
		}
		PDVS = pdvs
	})
	return PDVS
}

// GetInstance return wrapper for sql.DB
func GetInstance() *BeeruDatabase {
	databaseOnce.Do(func() {
		db, err := sql.Open("postgres", "user=caires dbname=postgres sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}
		Database = &BeeruDatabase{db}
	})
	return Database
}

// LoadToDatabase copy static data to database. For init proposal
func (b *BeeruDatabase) LoadToDatabase() {
	PDVS := getPDVS()
	for _, pdv := range PDVS.PDVS {
		if len(pdv.Document) != 17 {
			continue
		}
		bytesCoverageArea, _ := json.Marshal(pdv.CoverageArea)
		bytesAddress, _ := json.Marshal(pdv.Address)
		query := fmt.Sprintf("INSERT INTO pdv(id, trading_name,owner_name,document,coverage_area, address) VALUES($1,$2,$3,$4, ST_SetSRID(ST_GeomFromGeoJSON('%s'), 4326),ST_SetSRID(ST_GeomFromGeoJSON('%s'), 4326))", string(bytesCoverageArea), string(bytesAddress))
		_, err := b.DB.Exec(query, pdv.ID, pdv.TradingName, pdv.OwnerName, pdv.Document)
		if err != nil {
			log.Println(err, pdv.ID)
		}
	}
}
