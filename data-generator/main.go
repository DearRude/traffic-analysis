package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/planar"

	"google.golang.org/protobuf/proto"

	pb "github.com/DearRude/traffic-analysis/data-generator/protos/traffic" // protobuf
)

import _ "embed"

type TileName struct {
	Tile maptile.Tile
	Name string
}

type Point struct {
	Lat float64
	Lon float64
}

type LineTraffic struct {
	ID         uint32
	Length     float64
	Timestamp  time.Time
	City       string
	RoadClass  string
	Congestion string
	Geometry   []Point
}

var (
	//go:embed cities.geojson
	geojsonData []byte

	ZOOM maptile.Zoom
)

func main() {
	c := GenConfig()
	ZOOM = maptile.Zoom(c.Zoom)

	features, err := readGeoJSON()
	if err != nil {
		panic(err)
	}

	tileNames, err := genTileNames(features)
	if err != nil {
		panic(err)
	}

	for _, tileName := range tileNames {
		traffics, err := getTraffic(tileName, c.TileURL)
		if err != nil {
			log.Printf("error: %v", err) // Print the error
			continue
		}
		log.Printf("info: %d objects in tile %d %d scraped in %s", len(traffics.Traffics), tileName.Tile.X, tileName.Tile.Y, tileName.Name)

		_, err = proto.Marshal(traffics)
		if err != nil {
			fmt.Println("error encoding LineTraffics: ", err)
		}
	}

}

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

func getTraffic(tn TileName, url string) (*pb.LineTraffics, error) {
	bytes, err := requestTraffic(tn.Tile.X, tn.Tile.Y, url)
	if err != nil {
		return nil, fmt.Errorf("error request traffic data: %w", err)
	}

	ters, err := processTile(bytes, tn)
	if err != nil {
		return nil, fmt.Errorf("error process tile binary: %w", err)
	}

	return &pb.LineTraffics{Traffics: ters}, nil
}

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

func processTile(tile []byte, tn TileName) ([]*pb.LineTraffic, error) {
	var traffics []*pb.LineTraffic
	now := time.Now().UTC().Unix()

	tileData, err := mvt.UnmarshalGzipped(tile)
	if err != nil {
		return nil, fmt.Errorf("error unmarshal tile data: %w", err)
	}

	tileData.ProjectToWGS84(maptile.New(tn.Tile.X, tn.Tile.Y, ZOOM))
	features := tileData.ToFeatureCollections()["traffic"].Features

	for _, feature := range features {
		linestring, ok := feature.Geometry.(orb.LineString)
		if !ok {
			multiLineString, ok := feature.Geometry.(orb.MultiLineString)
			if !ok {
				return nil, fmt.Errorf("error unmarshal linestring or multilinestring")
			}

			// Merge the MultiLineString into a single LineString
			linestring = mergeMultiLineString(multiLineString)
		}
		feature.ID = generateLinestringID(linestring)

		traffic := &pb.LineTraffic{
			Id:         generateLinestringID(linestring),
			Length:     planar.Length(feature.Geometry),
			Timestamp:  now,
			City:       tn.Name,
			RoadClass:  feature.Properties.MustString("road_class", "unknown"),
			Congestion: feature.Properties.MustString("congestion", "unknown"),
			Geometry:   convertLineStringToPoints(linestring),
		}
		traffics = append(traffics, traffic)
	}

	return traffics, nil
}

func convertLineStringToPoints(lineString orb.LineString) []*pb.Point {
	points := make([]*pb.Point, len(lineString))
	for i, pt := range lineString {
		points[i] = &pb.Point{Lon: pt[0], Lat: pt[1]}
	}
	return points
}
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
