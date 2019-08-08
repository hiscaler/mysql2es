MySQL2ES
=========
Sync MySQL to ElasticSearch

### 编译
```bash
go build main.go m2e
```

### 运行
```bash
./m2e
```

## 设置
配置文件 `config/config.json`
```json
{
  "debug": true,
  "db": {
    "host": "127.0.0.1",
    "port": "3306",
    "databaseName": "database-name",
    "username": "root",
    "password": "root"
  },
  "dbOptions": {
    "syncTables": [
      "*"
    ],
    "ignoreTables": [],
    "mergeTables": {
      "www_merge_all": [
        "www_merge_1",
        "www_merge_2"
      ]
    },
    "defaultPk": "id",
    "datetimeFormatFields": [
      "created_at",
      "updated_at"
    ],
    "tables": {
      "www_category": {
        "datetimeFormatFields": [
          "created_at",
          "updated_at"
        ]
      }
    }
  },
  "es": {
    "indexPrefix": "",
    "urls": [
      "http://127.0.0.1:9200"
    ],
    "baseAuth": {
      "username": "",
      "password": ""
    }
  },
  "sizePerTime": 100
}
```
