package main

import (
	"flag"
	"log"
	"os"

	"github.com/peterbourgon/ff/v3"
)

type Config struct {
	// Requireds
	TileURL string

	// Optionals
	Zoom uint
}

func GenConfig() Config {
	log.Println("Read configurations.")
	fs := flag.NewFlagSet("data-generator", flag.ContinueOnError)
	var (
		tileURL = fs.String("tileURL", "", "Traffic tile url. e.g. https://example-tile.com/data/v2/%d/%d/%d.pbf")
		zoom    = fs.Uint("zoom", 14, "zoom level of tile")
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
		TileURL: *tileURL,
		Zoom:    *zoom,
	}
}
