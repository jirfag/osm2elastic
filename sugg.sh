curl -s -XPOST 'localhost:9200/osm/_suggest?pretty&size=20' -d '{
	"text" : "Ас",
	"osm-suggest" : {
		"completion" : {
			"field" : "SuggestData"
		}
	}
}'
