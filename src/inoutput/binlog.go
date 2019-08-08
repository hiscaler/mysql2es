package inoutput

type Binlog struct {
}

func (b Binlog) Read() ([]map[string]interface{}, error) {
	items := make([]map[string]interface{}, 0)
	
	return items, nil
}

func (b Binlog) Write() error {
	return nil
}
