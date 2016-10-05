package main

import (
	"flag"
	"log"
	"osm2elastic/osm"
	"osm2elastic/elastic"
)

func main() {
	var osmFilePath, elasticAddr string
	flag.StringVar(&osmFilePath, "osm-file", "planet.osm", "Path to .osm file")
	flag.StringVar(&elasticAddr, "elastic-addr", "127.0.0.1:9200", "Address of elasticsearch")
	flag.Parse()

	log.Printf("Parsing file %q", osmFilePath)
	nodes := osm.DecodeOsmNodes(osmFilePath)
	log.Printf("Parsed %d osm nodes from file %q", len(nodes), osmFilePath)

	elastic.ElasticImportOsmNodes(elasticAddr, nodes)
}
