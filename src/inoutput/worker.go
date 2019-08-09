package inoutput

type Worker interface {
	// Init
	Init() (err error)
	// Read data from MySQL database
	Read() (err error)
	// Write data to ElasticSearch
	Write() (insertCount, updateCount, deleteCount int, err error)
}

type ESItem struct {
	TableName string
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
