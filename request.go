package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/planar"
)

// Generate all tiles that need to be requested
func genTileNames(features []*geojson.Feature) ([]TileName, error) {
	var tileNames []TileName

	for _, feature := range features {
		polygon, ok := feature.Geometry.(orb.Polygon)
		if !ok {
			return nil, fmt.Errorf("error geojson feature is not polygon")
		}
		minTile := maptile.At(polygon.Bound().Min, ZOOM)
		maxTile := maptile.At(polygon.Bound().Max, ZOOM)
		for x := minTile.X; x <= maxTile.X; x++ {
			for y := maxTile.Y; y <= minTile.Y; y++ {
				tile := maptile.New(x, y, ZOOM)
				if !tileWithinPolygon(tile, polygon) {
					continue
				}
				tileName := TileName{
					Name: feature.Properties.MustString("name", " "),
					Tile: tile}
				tileNames = append(tileNames, tileName)
			}
		}
	}
	return tileNames, nil
}

// ReadGeoJSON reads the GeoJSON data from a embded binary
func readGeoJSON() ([]*geojson.Feature, error) {
	fc, err := geojson.UnmarshalFeatureCollection(geojsonData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal geojson file: %w", err)
	}
	return fc.Features, nil
}

// Function to check if a tile is within a polygon
func tileWithinPolygon(tile maptile.Tile, polygon orb.Polygon) bool {
	tileBound := tile.Bound()
	corners := []orb.Point{
		tileBound.Min,
		{tileBound.Max[0], tileBound.Min[1]},
		tileBound.Max,
		{tileBound.Min[0], tileBound.Max[1]},
	}

	for _, corner := range corners {
		if planar.PolygonContains(polygon, corner) {
			return true
		}
	}
	return false
}

// Requests and aggrigate a single tile
func getTraffic(tn TileName, url string) error {
	bytes, err := requestTraffic(tn.Tile.X, tn.Tile.Y, url)
	if err != nil {
		return fmt.Errorf("error request traffic data: %w", err)
	}

	if err := processTile(bytes, tn); err != nil {
		return fmt.Errorf("error process tile binary: %w", err)
	}

	return nil
}

// Does an API request to traffic tile provider
func requestTraffic(x, y uint32, url string) ([]byte, error) {
	url = fmt.Sprintf(url, ZOOM, x, y)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers to mimic a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, fmt.Errorf("no data available for the tile")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return body, nil
}
