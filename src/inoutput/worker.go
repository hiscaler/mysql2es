package inoutput

type Worker interface {
	// Read data from MySQL database
	Read() error
	// Write data to ElasticSearch
	Write() error
}

type ESItem struct {
	IndexName string
	IdName    string
	IdValue   string
	Values    map[string]interface{}
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
