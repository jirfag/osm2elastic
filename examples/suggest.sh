curl -s -XPOST 'localhost:9200/osm/_suggest?pretty&size=20' -d '{
	"text": "Актау",
	"osm-suggest": {
		"completion": {
			"field": "SuggestData"
		}
	}
}'
