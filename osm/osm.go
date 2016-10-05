package osm

import (
	"encoding/xml"
	"log"
	"os"
	"time"
)

// Location struct
type Location struct {
	Type        string
	Coordinates []float64
}

type Tag struct {
	XMLName xml.Name `xml:"tag"`
	Key     string   `xml:"k,attr"`
	Value   string   `xml:"v,attr"`
}

// Elem is a OSM base element
type Elem struct {
	ID        int64 `xml:"id,attr"`
	Loc       Location
	Version   int       `xml:"version,attr"`
	Ts        time.Time `xml:"timestamp,attr"`
	UID       int64     `xml:"uid,attr"`
	User      string    `xml:"user,attr"`
	ChangeSet int64     `xml:"changeset,attr"`
}

// Node structure
type Node struct {
	Elem
	XMLName xml.Name `xml:"node"`
	Lat     float64  `xml:"lat,attr"`
	Lng     float64  `xml:"lon,attr"`
	Tag     []Tag    `xml:"tag"`
}

func DecodeOsmNodes(fileName string) []Node {
	nodes := []Node{}
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Can't open OSM file: " + err.Error())
	}

	decoder := xml.NewDecoder(file)
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}

		switch typedToken := token.(type) {
		case xml.StartElement:
			if typedToken.Name.Local == "way" || typedToken.Name.Local == "relation" {
				log.Fatalln("Found not-node element in OSM file: " + typedToken.Name.Local)
			}
			if typedToken.Name.Local != "node" {
				continue
			}

			var n Node
			err = decoder.DecodeElement(&n, &typedToken)
			if err != nil {
				log.Fatalln(err)
			}
			nodes = append(nodes, n)
		}
	}
	return nodes
}
