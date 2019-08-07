package main

import (
	"config"
	"fmt"
	"github.com/go-ozzo/ozzo-dbx"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var cfg *config.Config
var db *dbx.DB

func init() {
	cfg = config.NewConfig()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	db, err = dbx.Open("mysql", cfg.DB.Username+":"+cfg.DB.Password+"@/"+cfg.DB.DatabaseName)
	if err != nil {
		log.Panicln("Open database error: " + err.Error())
	}
}

func main() {
	row := dbx.NullStringMap{}
	tables := make([]string, 0)
	db.NewQuery("SHOW TABLES").Column(&tables)
	for _, table := range tables {
		fmt.Println("Table: " + table)
		fmt.Println("SELECT * FROM " + table + " LIMIT 10")
		sq := db.Select().From(table).Limit(10)
		rows, err := sq.Rows()
		if err == nil {
			for rows.Next() {
				rows.ScanMap(row)
				item := make(map[string]interface{})
				for fieldName, v := range row {
					fieldValue, _ := v.Value()
					item[fieldName] = fieldValue
				}
				fmt.Println(fmt.Sprintf("#%v", item))
			}
		}
	}
}
