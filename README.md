MySQL2ES
=========
Sync MySQL to ElasticSearch

## 设置
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
    "defaultPk": "id",
    "datetimeFormatFields": [
      "created_at",
      "updated_at"
    ],
    "ignoreTables": [],
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
