package inoutput

type Worker interface {
	// Read data from MySQL database
	Read() ([]map[string]interface{}, error)
	// Write data to ElasticSearch
	Write() error
}

func In(value string, items []string) bool {
	exists := false
	for _, v := range items {
		if value == v {
			exists = true
			break
		}
	}
	return exists
}
