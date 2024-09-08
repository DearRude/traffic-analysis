package main

import (
	"fmt"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"github.com/paulmach/orb/maptile"
)

type TileName struct {
	Tile maptile.Tile
	Name string
}

type Way struct {
	WayID     uint32 `gorm:"primary_key"`
	Length    float64
	City      string
	RoadClass string
	Geom      []byte `gorm:"type:geometry(LineString,4326)"`
}

type Traffic struct {
	ID           uint `gorm:"primary_key"`
	WayID        uint32
	TrafficState string
	RecordTime   time.Time
}

// Insert ways using correct geom to PostGIS
func insertWays(ways []Way, lines []orb.LineString) error {
	// Start a transaction
	tx := dbase.Begin()

	// Prepare SQL statement and values
	sql := "INSERT INTO ways (way_id, length, city, road_class, geom) VALUES "
	var values []interface{}

	for i, way := range ways {
		sql += "(?, ?, ?, ?, ST_GeomFromEWKB(?)),"
		values = append(values, way.WayID, way.Length, way.City, way.RoadClass, ewkb.Value(lines[i], 4326))
	}
	sql = sql[:len(sql)-1] // Remove trailing comma
	sql += " ON CONFLICT (way_id) DO NOTHING"

	// Execute SQL
	if err := tx.Exec(sql, values...).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert ways: %w", err)
	}

	// Commit transaction
	tx.Commit()
	return nil
}
