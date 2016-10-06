package elastic

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/jirfag/osm2elastic/osm"
	"gopkg.in/olivere/elastic.v2"
)

const (
	elasticIndexName = "osm"
	elasticTypeName = "node"
)

func ElasticImportOsmNodes(elasticAddr string, nodes []osm.Node) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(fmt.Sprintf("http://%s", elasticAddr)))
	if err != nil {
		panic(err)
	}

	client.DeleteIndex(elasticIndexName).Do()

	// Create an index
	_, err = client.CreateIndex(elasticIndexName).Do()
	if err != nil {
		panic(err)
	}

	CreateMapping(client)

	nodeBatches := GroupNodesToBatches(nodes, 30)
	ImportNodeBatches(client, nodeBatches, len(nodes))
}

func CreateMapping(client *elastic.Client) {
	mappingBody, err := ioutil.ReadFile("es_config/mappings.json")
	if err != nil {
		log.Fatalf("can't read mappings file: %s", err)
	}

	putresp, err := client.PutMapping().Index(elasticIndexName).Type(elasticTypeName).BodyString(string(mappingBody)).Do()
	if err != nil {
		log.Fatalf("expected put mapping to succeed; got: %v", err)
	}
	if putresp == nil {
		log.Fatalf("expected put mapping response; got: %v", putresp)
	}
	if !putresp.Acknowledged {
		log.Fatalf("expected put mapping ack; got: %v", putresp.Acknowledged)
	}
}

type NodeBatch struct {
	nodes []osm.Node
}

func GroupNodesToBatches(nodes []osm.Node, batchSize int) []NodeBatch {
	var batches []NodeBatch
	currentBatch := NodeBatch{}
	for _, node := range nodes {
		currentBatch.nodes = append(currentBatch.nodes, node)
		if len(currentBatch.nodes) == batchSize {
			batches = append(batches, currentBatch)
			currentBatch = NodeBatch{}
		}
	}
	if len(currentBatch.nodes) == batchSize {
		batches = append(batches, currentBatch)
	}

	return batches
}

func ImportNodeBatches(client *elastic.Client, nodeBatches []NodeBatch, totalNodes int) {
	importedNodesCount := 0
	nextLogPrintTime := time.Now()
	for _, nodeBatch := range nodeBatches {
		ImportNodeBatch(client, nodeBatch)
		importedNodesCount += len(nodeBatch.nodes)
		if time.Now().After(nextLogPrintTime) {
			nextLogPrintTime = time.Now().Add(time.Second * 5)
			log.Printf("imported %d/%d nodes", importedNodesCount, totalNodes)
		}
	}
}

func ImportNodeBatch(client *elastic.Client, nodeBatch NodeBatch) {
	bulkRequest := client.Bulk()
	for _, node := range nodeBatch.nodes {
		doc := NodeToDoc(node)
		if doc == nil { // skip some bad nodes
			continue
		}
		IdStr := strconv.FormatInt(node.ID, 10)
		indexReq := elastic.NewBulkIndexRequest().Index(elasticIndexName).Type(elasticTypeName).Id(IdStr).Doc(doc)
		bulkRequest = bulkRequest.Add(indexReq)
	}
	bulkResponse, err := bulkRequest.Do()
	if err != nil {
		log.Fatalln(err)
	}
	if bulkResponse == nil {
		log.Fatalln("no response")
	}
	if bulkResponse.Errors {
		log.Fatalln("got errors")
	}
}

type NodeDoc struct {
	SuggestData NodeDocSuggestData
}

type NodeDocInfo struct {
	OsmId int64
	Name string // default name; in some cases there is only 'name'
	NameRu string // russian
	NameKk string // kazakhstan
	NameEn string // english
	Population int // how many people
	Lat float64 // lattitude
	Lon float64 // longitude

	Country string
	Region string
	PlaceType string // city OR town OR village OR hamlet 
}

type NodeDocSuggestData struct {
	Input []string `json:"input"`
	Output string `json:"output"`
	Payload NodeDocInfo `json:"payload"`
	Weight int `json:"weight"`
}

type NodeDocSuggestDataPayload struct {
	Id int `json:"id"`
}

func NodeToDoc(n osm.Node) *NodeDoc {
	nodeTags := map[string]string{}
	for _, tag := range n.Tag {
		nodeTags[tag.Key] = tag.Value
	}

	population, _ := strconv.Atoi(nodeTags["population"])

	nodeDocInfo := NodeDocInfo{
		n.ID, // Id
		nodeTags["name"], // Name
		nodeTags["name:ru"], // NameRu
		nodeTags["name:kk"], // NameKk
		nodeTags["name:en"], // NameEn
		population, // Population
		n.Lat, // Lat
		n.Lng, // Lon
		nodeTags["addr:country"], // Country
		nodeTags["addr:region"], // Region
		nodeTags["place"], // PlaceType
	}
	suggestData := NodeDocSuggestData{
		[]string{nodeDocInfo.Name, nodeDocInfo.NameRu, nodeDocInfo.NameKk, nodeDocInfo.NameEn}, // Input
		strconv.FormatInt(n.ID, 10), // Output
		nodeDocInfo, // Payload
		nodeDocInfo.Population, // Weight
	}
	ret := NodeDoc{suggestData}

	// Primary languages should be set
	if nodeDocInfo.Name == "" && nodeDocInfo.NameRu == "" && nodeDocInfo.NameKk == "" {
		log.Printf("skip node: no primary languages for node %v", n)
		return nil
	}

	return &ret
}
