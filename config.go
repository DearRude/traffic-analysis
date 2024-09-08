package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/peterbourgon/ff/v3"
)

type Config struct {
	// Requireds
	TileURL      string
	PostGisParam string

	// Optionals
	Zoom           uint
	MetricsAddr    string
	ScrapeInterval time.Duration
}

func GenConfig() Config {
	log.Println("Read configurations.")
	fs := flag.NewFlagSet("traffic-analysis", flag.ContinueOnError)
	var (
		tileUrl        = fs.String("tileUrl", "", "Traffic tile url. e.g. https://example-tile.com/data/v2/%d/%d/%d.pbf")
		postGisParam   = fs.String("postGisParam", "host=localhost user=admin password=trafficramz dbname=postgres port=5433 sslmode=disable TimeZone=UTC", "postgres/postgis url connection")
		zoom           = fs.Uint("zoom", 14, "zoom level of tile")
		metricsAddr    = fs.String("metricsAddr", ":8080", "address ip/port of metric webserver")
		scrapeInterval = fs.Duration("scrapeInterval", 10*time.Minute, "interval of each scrape to happen")
	)

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fs.String("config", "", "config file")
	} else {
		fs.String("config", ".env", "config file")
	}

	err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVars(),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.EnvParser),
	)
	if err != nil {
		log.Fatalf("Unable to parse args. Error: %s", err)
	}

	return Config{
		TileURL:      *tileUrl,
		PostGisParam: *postGisParam,

		Zoom:           *zoom,
		MetricsAddr:    *metricsAddr,
		ScrapeInterval: *scrapeInterval,
	}
}
