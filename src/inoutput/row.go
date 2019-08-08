package inoutput

import (
	"config"
	"context"
	"fmt"
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/olivere/elastic"
	"log"
	"strconv"
	"time"
)

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

	var options []elastic.ClientOptionFunc
	if len(cfg.ES.Urls) > 0 {
		options = append(options, elastic.SetURL(cfg.ES.Urls...))
	}
	if len(cfg.ES.BaseAuth.Username) > 0 && len(cfg.ES.BaseAuth.Password) > 0 {
		options = append(options, elastic.SetBasicAuth(cfg.ES.BaseAuth.Username, cfg.ES.BaseAuth.Password))
	}
	esClient, err = elastic.NewClient(options...)
	if err != nil {
		log.Panicln(err)
	}

	ctx = context.Background()
}

type Row struct {
	TableIndexes map[string]string
	Items        []ESItem
}

func (r *Row) Init() (err error) {
	if r.TableIndexes == nil {
		r.TableIndexes = make(map[string]string, 0)
	}
	tables := make([]string, 0)
	dbOptions := cfg.DBOptions
	db.NewQuery("SHOW TABLES").Column(&tables)
	for _, table := range tables {
		if !In(table, dbOptions.IgnoreTables) {
			// 检测 ES index 是否存在
			indexName := table
			for k, v := range dbOptions.MergeTables {
				if In(table, v) {
					indexName = k
					break
				}
			}
			r.TableIndexes[table] = indexName
			exists := false
			exists, err = esClient.IndexExists(indexName).Do(ctx)
			if err != nil {
				log.Panicln(err)
			} else if !exists {
				log.Println(fmt.Sprintf("Create ES `%s` index", indexName))
				esClient.CreateIndex(indexName).Do(ctx)
			}
		}
	}

	return err
}

func (r *Row) Read() (err error) {
	dbOptions := cfg.DBOptions
	for table, indexName := range r.TableIndexes {
		row := dbx.NullStringMap{}
		pkName := ""
		pkValue := ""
		ignoreFields := make([]string, 0)
		datetimeFormatFields := dbOptions.DatetimeFormatFields
		for k, v := range dbOptions.Tables {
			if k == table {
				if len(v.PK) == 0 {
					pkName = dbOptions.DefaultPk
				}
				ignoreFields = v.IgnoreFields
				datetimeFormatFields = append(datetimeFormatFields, v.DatetimeFormatFields...)
				break
			}
		}
		if len(pkName) == 0 {
			pkName = dbOptions.DefaultPk
		}
		sq := db.Select().From(table).Limit(cfg.SizePerTime)
		var rows *dbx.Rows
		rows, err = sq.Rows()
		if err == nil {
			for rows.Next() {
				rows.ScanMap(row)
				item := ESItem{
					IndexName: indexName,
					IdName:    pkName,
				}
				values := make(map[string]interface{})
				for fieldName, v := range row {
					if In(fieldName, ignoreFields) {
						continue
					}
					fieldValue, _ := v.Value()
					if fieldName == pkName {
						pkValue = v.String
					}
					item.IdValue = pkValue
					if In(fieldName, datetimeFormatFields) {
						fieldName += "_formatted"
						v, _ := strconv.ParseInt(fieldValue.(string), 10, 64)
						fieldValue = time.Unix(v, 0)
					}
					values[fieldName] = fieldValue
				}
				item.Values = values
				r.Items = append(r.Items, item)
			}
		}
	}

	return
}

func (r *Row) Write() (insertCount, updateCount, deleteCount int, err error) {
	var e error
	for _, item := range r.Items {
		q, err := esClient.Search(item.IndexName).Query(elastic.NewTermQuery(item.IdName, item.IdValue)).Do(ctx)
		if err == nil {
			if q.TotalHits() == 0 {
				put, err := esClient.Index().
					Index(item.IndexName).
					Id(item.IdValue).
					BodyJson(item.Values).
					Do(ctx)
				if err != nil {
					log.Printf("IndexName: %s, IdName: %s, IdValue: %s, err: %v", item.IndexName, item.IdName, item.IdValue, err)
				}
				insertCount++
				log.Printf("Indexed `%s` to `%s` index, type `%s`\n", put.Id, put.Index, put.Type)
			} else {
				put, err := esClient.Update().
					Index(item.IndexName).
					Id(item.IdValue).
					Doc(item.Values).
					Do(ctx)
				if err != nil {
					log.Panicf("IndexName: %s, IdName: %s, IdValue: %s, err: %v", item.IndexName, item.IdName, item.IdValue, err)
				}
				updateCount++
				log.Printf("Update `%s` to `%s` index, type `%s`\n", put.Id, put.Index, put.Type)
			}
		} else {
			e = err
		}
	}

	return insertCount, updateCount, deleteCount, e
}
