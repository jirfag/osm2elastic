curl -s -XPOST 'localhost:9200/osm/_search?pretty&size=20' -d '{
	"query": {
		"function_score": {
			"functions": [
				{
					"linear": {
						"Location": { 
							"origin": { "lat": 55.7972075, "lon": 37.5355795 },
							"offset": "10km",
							"scale":  "100000km"
						}
					},
					"weight": 1.0
				},
				{
					"linear": {
						"Population": { 
							"origin": "100000000", 
							"offset": "1000",
							"scale":  "100000000"
						}
					},
					"weight": 1.5
				}
			],
			"score_mode": "avg",
			"boost_mode": "replace"
		}
	}
}'
