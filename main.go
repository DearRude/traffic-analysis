package main

import (
	"fmt"
	"hash/crc32"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/planar"
)

import _ "embed"

var (
	//go:embed cities.geojson
	geojsonData []byte
	dbase       *gorm.DB

	ZOOM maptile.Zoom
)

func main() {
	c := GenConfig()
	ZOOM = maptile.Zoom(c.Zoom)

	// init database
	db, err := gorm.Open(postgres.Open(c.PostGisParam), &gorm.Config{CreateBatchSize: 500})
	if err != nil {
		panic(err)
	}
	dbase = db
	if err := dbase.AutoMigrate(&Way{}, &Traffic{}); err != nil {
		panic(err)
	}
	log.Println("Auto migrated to new schemas")

	// read features from geojson
	features, err := readGeoJSON()
	if err != nil {
		panic(err)
	}

	// generate tile names
	tileNames, err := genTileNames(features)
	if err != nil {
		panic(err)
	}

	// init metrics
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	go func() {
		log.Fatal(http.ListenAndServe(c.MetricsAddr, nil))
	}()

	ticker := time.NewTicker(c.ScrapeInterval)
	for range ticker.C {
		resetCongestionMetrics()
		// sem is a channel that will allow up to 10 concurrent operations.
		var sem = make(chan int, 10)
		for _, tileName := range tileNames {
			tileName := tileName
			sem <- 1
			go func() {
				err := getTraffic(tileName, c.TileURL)
				<-sem

				if err != nil {
					log.Printf("error: %v", err) // Print the error
					return
				}
				log.Printf("info: objects in tile %d %d scraped in %s", tileName.Tile.X, tileName.Tile.Y, tileName.Name)
			}()
		}
		setCongestionMetrics(m)
	}
}

// Aggrigate each tile and save them to db
func processTile(tile []byte, tn TileName) error {
	var ways []Way
	var traffics []Traffic
	var lines []orb.LineString

	now := time.Now().UTC()
	tileData, err := mvt.UnmarshalGzipped(tile)
	if err != nil {
		return fmt.Errorf("error unmarshal tile data: %w", err)
	}

	tileData.ProjectToWGS84(maptile.New(tn.Tile.X, tn.Tile.Y, ZOOM))
	features := tileData.ToFeatureCollections()["traffic"].Features

	for _, feature := range features {
		linestring, ok := feature.Geometry.(orb.LineString)
		if !ok {
			multiLineString, ok := feature.Geometry.(orb.MultiLineString)
			if !ok {
				return fmt.Errorf("error unmarshal linestring or multilinestring")
			}

			// Merge the MultiLineString into a single LineString
			linestring = mergeMultiLineString(multiLineString)
		}
		feature.ID = generateLinestringID(linestring)
		uniqueWayId := generateLinestringID(linestring)

		lines = append(lines, linestring)

		way := Way{
			WayID:     uniqueWayId,
			Length:    planar.Length(feature.Geometry) * 10000, // planar distance does not account for earth curve
			City:      tn.Name,
			RoadClass: feature.Properties.MustString("road_class", "unknown"),
		}
		ways = append(ways, way)

		traffic := Traffic{
			WayID:        uniqueWayId,
			TrafficState: feature.Properties.MustString("congestion", "unknown"),
			RecordTime:   now,
		}
		traffics = append(traffics, traffic)
	}

	writeCongestionMetrics(ways, traffics)

	// Add the data to postgres
	if err := insertWays(ways, lines); err != nil {
		return fmt.Errorf("error inserting ways to db: %w", err)
	}
	dbase.Table("traffics").Clauses(clause.OnConflict{DoNothing: true}).Create(&traffics)

	return nil
}

// Generate unique id for each linestring based on their points
func generateLinestringID(points []orb.Point) uint32 {
	var normalized string
	for _, point := range points {
		normalized += fmt.Sprintf("%f,%f;", point.X(), point.Y())
	}
	return crc32.ChecksumIEEE([]byte(normalized))
}

// Helper function to merge a MultiLineString into a single LineString
func mergeMultiLineString(mls orb.MultiLineString) orb.LineString {
	var merged orb.LineString
	for _, line := range mls {
		merged = append(merged, line...)
	}

	return merged
}
