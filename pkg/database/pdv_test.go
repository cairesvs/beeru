package database

import (
	"strconv"
	"testing"

	"github.com/cairesvs/beeru/pkg/model"
)

var createTests = []struct {
	about      string
	pdvs       []*model.PDV
	shouldFail bool
}{
	{
		about:      "Create simple PDV",
		shouldFail: false,
		pdvs: []*model.PDV{
			{
				ID:          "1",
				TradingName: "Ronaldo's",
				OwnerName:   "Ronaldo",
				Document:    "02.453.716/000170",
				CoverageArea: &model.MultiPolygon{
					Type: "MultiPolygon",
					Coordinates: [][][][]float64{
						{
							{
								{

									-46.65393590927124,
									-23.56369210202025,
								},
								{

									-46.656789779663086,
									-23.571008561959243,
								},
								{

									-46.649043560028076,
									-23.57437163664494,
								},
								{

									-46.64374351501465,
									-23.56912048228182,
								},
								{

									-46.64666175842285,
									-23.564793530648632,
								},
								{

									-46.65393590927124,
									-23.56369210202025,
								},
							},
						},
					},
				},
				Address: &model.Point{
					Type: "Point",
					Coordinates: []float64{
						-46.647348403930664,
						-23.570615214264883,
					},
				},
			},
		},
	}, {
		about:      "Couldn't create the with same document filed",
		shouldFail: true,
		pdvs: []*model.PDV{
			{
				ID:          "1",
				TradingName: "Ronaldo's",
				OwnerName:   "Ronaldo",
				Document:    "02.453.716/000170",
				CoverageArea: &model.MultiPolygon{
					Type: "MultiPolygon",
					Coordinates: [][][][]float64{
						{
							{
								{

									-46.65393590927124,
									-23.56369210202025,
								},
								{

									-46.656789779663086,
									-23.571008561959243,
								},
								{

									-46.649043560028076,
									-23.57437163664494,
								},
								{

									-46.64374351501465,
									-23.56912048228182,
								},
								{

									-46.64666175842285,
									-23.564793530648632,
								},
								{

									-46.65393590927124,
									-23.56369210202025,
								},
							},
						},
					},
				},
				Address: &model.Point{
					Type: "Point",
					Coordinates: []float64{
						-46.647348403930664,
						-23.570615214264883,
					},
				},
			}, {
				ID:          "2",
				TradingName: "Romario's",
				OwnerName:   "Romario",
				Document:    "02.453.716/000170",
				CoverageArea: &model.MultiPolygon{
					Type: "MultiPolygon",
					Coordinates: [][][][]float64{
						{
							{
								{

									-46.65393590927124,
									-23.56369210202025,
								},
								{

									-46.656789779663086,
									-23.571008561959243,
								},
								{

									-46.649043560028076,
									-23.57437163664494,
								},
								{

									-46.64374351501465,
									-23.56912048228182,
								},
								{

									-46.64666175842285,
									-23.564793530648632,
								},
								{

									-46.65393590927124,
									-23.56369210202025,
								},
							},
						},
					},
				},
				Address: &model.Point{
					Type: "Point",
					Coordinates: []float64{
						-46.647348403930664,
						-23.570615214264883,
					},
				},
			},
		},
	}}

func TestCreateTest(t *testing.T) {
	beeru := GetInstance()
	for _, test := range createTests {
		beeru.Truncate()
		t.Log(test.about)
		ids := []int64{}
		for _, p := range test.pdvs {
			input := &PDVCreateInput{
				Database: beeru,
				PDV:      p,
			}
			id, err := CreatePDV(input)
			if err != nil && !test.shouldFail {
				t.Errorf("Couldn't create PDV %s", err)
			}
			if id != 0 {
				ids = append(ids, id)
			}
		}
		if len(ids) == 0 {
			t.Error("No ID returned")
		}
		for i, p := range ids {
			input := &PDVGetInput{
				Database: beeru,
				ID:       strconv.FormatInt(p, 10),
			}
			pdv := GetPDV(input)
			if pdv == nil {
				t.Error("PDV wasn't inserted")
			}
			if pdv.OwnerName != test.pdvs[i].OwnerName {
				t.Errorf("PDV OwnerName doesn't match %s %s", pdv.OwnerName, test.pdvs[i].OwnerName)
			}
			if pdv.TradingName != test.pdvs[i].TradingName {
				t.Errorf("PDV TradingName doesn't match %s %s", pdv.TradingName, test.pdvs[i].TradingName)
			}
			if pdv.Document != test.pdvs[i].Document {
				t.Errorf("PDV Document doesn't match %s %s", pdv.Document, test.pdvs[i].Document)
			}
		}
	}
}

var findTests = []struct {
	about       string
	pdvs        []*model.PDV
	responsePDV []*model.PDV
	point       *model.Point
}{
	{
		about: "Find PDV for given point",
		point: &model.Point{
			Type: "Point",
			Coordinates: []float64{
				-46.647348403930664,
				-23.570615214264883,
			},
		},
		responsePDV: []*model.PDV{
			{
				ID:          "1",
				TradingName: "Ronaldo's",
				OwnerName:   "Ronaldo",
				Document:    "02.453.716/000170",
				CoverageArea: &model.MultiPolygon{
					Type: "MultiPolygon",
					Coordinates: [][][][]float64{
						{
							{
								{

									-46.65393590927124,
									-23.56369210202025,
								},
								{

									-46.656789779663086,
									-23.571008561959243,
								},
								{

									-46.649043560028076,
									-23.57437163664494,
								},
								{

									-46.64374351501465,
									-23.56912048228182,
								},
								{

									-46.64666175842285,
									-23.564793530648632,
								},
								{

									-46.65393590927124,
									-23.56369210202025,
								},
							},
						},
					},
				},
				Address: &model.Point{
					Type: "Point",
					Coordinates: []float64{
						-46.647348403930664,
						-23.570615214264883,
					},
				},
			},
		},
		pdvs: []*model.PDV{
			{
				ID:          "1",
				TradingName: "Ronaldo's",
				OwnerName:   "Ronaldo",
				Document:    "02.453.716/000170",
				CoverageArea: &model.MultiPolygon{
					Type: "MultiPolygon",
					Coordinates: [][][][]float64{
						{
							{
								{

									-46.65393590927124,
									-23.56369210202025,
								},
								{

									-46.656789779663086,
									-23.571008561959243,
								},
								{

									-46.649043560028076,
									-23.57437163664494,
								},
								{

									-46.64374351501465,
									-23.56912048228182,
								},
								{

									-46.64666175842285,
									-23.564793530648632,
								},
								{

									-46.65393590927124,
									-23.56369210202025,
								},
							},
						},
					},
				},
				Address: &model.Point{
					Type: "Point",
					Coordinates: []float64{
						-46.647348403930664,
						-23.570615214264883,
					},
				},
			},
		},
	}, {
		about: "Should return empty results for point outside of coverage area",
		point: &model.Point{
			Type: "Point",
			Coordinates: []float64{
				-40.647348403930664,
				-23.570615214264883,
			},
		},
		responsePDV: []*model.PDV{},
		pdvs: []*model.PDV{
			{
				ID:          "1",
				TradingName: "Ronaldo's",
				OwnerName:   "Ronaldo",
				Document:    "02.453.716/000170",
				CoverageArea: &model.MultiPolygon{
					Type: "MultiPolygon",
					Coordinates: [][][][]float64{
						{
							{
								{

									-46.65393590927124,
									-23.56369210202025,
								},
								{

									-46.656789779663086,
									-23.571008561959243,
								},
								{

									-46.649043560028076,
									-23.57437163664494,
								},
								{

									-46.64374351501465,
									-23.56912048228182,
								},
								{

									-46.64666175842285,
									-23.564793530648632,
								},
								{

									-46.65393590927124,
									-23.56369210202025,
								},
							},
						},
					},
				},
				Address: &model.Point{
					Type: "Point",
					Coordinates: []float64{
						-46.647348403930664,
						-23.570615214264883,
					},
				},
			},
		},
	},
}

func TestFindTest(t *testing.T) {
	beeru := GetInstance()
	for _, test := range findTests {
		beeru.Truncate()
		t.Log(test.about)
		for _, p := range test.pdvs {
			input := &PDVCreateInput{
				Database: beeru,
				PDV:      p,
			}
			_, err := CreatePDV(input)
			if err != nil {
				t.Errorf("Couldn't create PDV %s", err)
			}
		}
		input := &PDVFindInput{
			Database: beeru,
			Point:    test.point,
		}
		pdvs := FindPDV(input)
		for i, pdv := range test.responsePDV {
			if pdv.OwnerName != pdvs.PDVS[i].OwnerName {
				t.Errorf("PDV OwnerName doesn't match %s %s", pdv.OwnerName, pdvs.PDVS[i].OwnerName)
			}
			if pdv.TradingName != pdvs.PDVS[i].TradingName {
				t.Errorf("PDV TradingName doesn't match %s %s", pdv.TradingName, pdvs.PDVS[i].TradingName)
			}
			if pdv.Document != pdvs.PDVS[i].Document {
				t.Errorf("PDV Document doesn't match %s %s", pdv.Document, pdvs.PDVS[i].Document)
			}
		}
	}
}
