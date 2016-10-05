package main

import (
	"flag"
	"log"
	"osm2elastic/osm"
	"osm2elastic/elastic"
)

func main() {
	var osmFilePath string
	flag.StringVar(&osmFilePath, "osm-file", "planet.osm", "Path to .osm file")
	flag.Parse()

	log.Printf("Parsing file %q", osmFilePath)
	nodes := osm.DecodeOsmNodes(osmFilePath)
	log.Printf("Parsed %d osm nodes from file %q", len(nodes), osmFilePath)

	elastic.ElasticImportOsmNodes(nodes)
}
