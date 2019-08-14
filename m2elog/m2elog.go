package m2elog

import (
	"fmt"
	"github.com/go-ozzo/ozzo-dbx"
	"github.com/hiscaler/mysql2es/config"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	TableName    = "m2e_log"
	PKIntType    = "int"
	PKStringType = "string"
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

type M2ELog struct {
	Id            int64
	TableName     string
	PkName        string
	PkStringValue string
	PkIntValue    int64
	LastModify    int64
	Version       string
}

func versionFields(eLog M2ELog) []string {
	versionFields := make([]string, 0)
	for k, v := range cfg.DBOptions.Tables {
		if k == eLog.TableName {
			if len(v.VersionFields) > 0 {
				versionFields = v.VersionFields
			} else {
				defaultVersionField := cfg.DBOptions.DefaultPk
				if len(v.PK) > 0 {
					defaultVersionField = v.PK
				}
				versionFields = append(versionFields, defaultVersionField)
			}
			break
		}
	}
	if len(versionFields) == 0 {
		versionFields = append(versionFields, cfg.DBOptions.DefaultPk)
	}

	return versionFields
}

func (eLog *M2ELog) PkType() (typ string) {
	typ = cfg.DBOptions.DefaultPkType
	for k, v := range cfg.DBOptions.Tables {
		if k == eLog.TableName {
			typ = v.PkType
			break
		}
	}

	if typ != PKIntType && typ != PKStringType {
		typ = PKIntType
	}

	return
}

// Insert OR Update log
func (eLog *M2ELog) Save() (isNewRecord, success bool, err error) {
	version := strings.TrimSpace(eLog.Version)
	if len(version) == 0 {
		versionFields := versionFields(*eLog)
		var row dbx.NullStringMap
		err := db.Select(versionFields...).From(eLog.TableName).Where(dbx.HashExp{eLog.PkName: eLog.PkIntValue}).One(&row)
		if err == nil {
			versions := make([]string, 0)
			for k, v := range row {
				value := ""
				if v.Valid {
					value = v.String
				}
				versions = append(versions, fmt.Sprintf("%s:%s", k, value))
			}
			version = strings.Join(versions, ",")
		}
	}
	var intValue int64
	var strValue string
	typ := eLog.PkType()
	if typ == PKIntType {
		intValue = eLog.PkIntValue
		if intValue == 0 {
			intValue, _ = strconv.ParseInt(eLog.PkStringValue, 10, 64)
		}
	} else if typ == PKStringType {
		strValue = eLog.PkStringValue
	}
	params := dbx.Params{
		"table_name":      eLog.TableName,
		"pk_name":         eLog.PkName,
		"pk_string_value": strValue,
		"pk_int_value":    intValue,
		"version":         version,
		"last_modify":     time.Now().Unix(),
	}
	var n int
	db.Select("COUNT(*)").From(TableName).Where(dbx.HashExp{"table_name": eLog.TableName, "pk_name": eLog.PkName, "pk_int_value": eLog.PkIntValue}).Row(&n)
	if n == 0 {
		isNewRecord = true
	}
	if isNewRecord {
		// Insert
		if _, err = db.Insert(TableName, params).Execute(); err == nil {
			success = true
		}
	} else {
		// Update
		if _, err = db.Update(TableName, params, dbx.HashExp{"id": eLog.Id}).Execute(); err == nil {
			success = true
		}
	}

	return
}

// Delete log
func (eLog *M2ELog) Delete() bool {
	if _, err := db.Delete(TableName, dbx.HashExp{"id": eLog.Id}).Execute(); err == nil {
		return true
	} else {
		return false
	}
}

// Check log status
func (eLog *M2ELog) Status() (changed, deleted bool) {
	versionFields := versionFields(*eLog)
	var row dbx.NullStringMap
	err := db.Select(versionFields...).From(eLog.TableName).Where(dbx.HashExp{eLog.PkName: eLog.PkIntValue}).One(&row)
	if err != nil {
		log.Println("Error: ", err)
	}
	if row == nil {
		deleted = true
	} else {
		versions := make(map[string]string)
		for _, v := range strings.Split(eLog.Version, ",") {
			if strings.Contains(v, ":") {
				item := strings.SplitN(v, ":", 2)
				versions[item[0]] = item[1]
			}
		}
		for k, v := range versions {
			for rowK, rowV := range row {
				if rowK == k {
					if rowV.Valid {
						if rowV.String != v {
							changed = true
							break
						}
					}
				}
				if changed {
					break
				}
			}
		}
	}

	return
}
