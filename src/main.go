package main

import (
	"config"
	"context"
	"fmt"
	"github.com/go-ozzo/ozzo-dbx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olivere/elastic"
	"log"
	"time"
)

const Default_PK_NAME = "id"

var cfg *config.Config
var db *dbx.DB
var esClient *elastic.Client
var ctx context.Context

func init() {
	cfg = config.NewConfig()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	db, err = dbx.Open("mysql", cfg.DB.Username+":"+cfg.DB.Password+"@/"+cfg.DB.DatabaseName)
	if err != nil {
		log.Panicln("Open database error: " + err.Error())
	}

	esClient, err = elastic.NewClient()
	if err != nil {
		log.Panicln(err)
	}

	ctx = context.Background()
}

func main() {
	insertRecords := 0
	updateRecords := 0
	deleteRecords := 0
	beginTimestamp := time.Now().Unix()
	fmt.Println("Begin sync")
	row := dbx.NullStringMap{}
	tables := make([]string, 0)
	db.NewQuery("SHOW TABLES").Column(&tables)
	pkName := Default_PK_NAME
	pkValue := ""
	for _, table := range tables {
		// 检测 ES index 是否存在
		exists, err := esClient.IndexExists(table).Do(ctx)
		if err != nil {
			panic(err)
		} else if !exists {
			fmt.Println(fmt.Sprintf("Create ES `%s` index", table))
			esClient.CreateIndex(table).Do(ctx)
		}

		fmt.Println("Table: " + table)
		sq := db.Select().From(table).Limit(10)
		rows, err := sq.Rows()
		if err == nil {
			indexService := esClient.Index().Index(table)
			for rows.Next() {
				rows.ScanMap(row)
				item := make(map[string]interface{})
				for fieldName, v := range row {
					fieldValue, _ := v.Value()
					if fieldName == pkName {
						pkValue = v.String
					}
					item[fieldName] = fieldValue
				}
				fmt.Println(fmt.Sprintf("#%v", item))
				q, err := esClient.Search(table).Query(elastic.NewTermQuery(pkName, pkValue)).Do(ctx)
				if err == nil {
					if q.TotalHits() == 0 {
						put, err := indexService.
							Id(pkValue).
							BodyJson(item).
							Do(ctx)
						if err != nil {
							panic(err)
						}
						insertRecords++
						log.Printf("Indexed trace %s to index %s, type %s\n", put.Id, put.Index, put.Type)
					} else {
						put, err := esClient.Update().
							Index(table).
							Id(pkValue).
							Doc(item).
							Do(ctx)
						if err != nil {
							panic(err)
						}
						updateRecords++
						log.Printf("Update trace %s to index %s, type %s\n", put.Id, put.Index, put.Type)
					}
				} else {
					panic(err)
				}
			}
		}
	}
	fmt.Println(fmt.Sprintf("Insert: %d, Update: %d, Delete: %d, cost %d seconds.", insertRecords, updateRecords, deleteRecords, time.Now().Unix()-beginTimestamp))
	fmt.Println("Done sync")
}
