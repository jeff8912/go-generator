package modules

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"{{.packagePrefix}}/config"
	"time"
)

var (
	DbInstance *gorm.DB

	SAVE_AFFECTED_ZERO_ERROR   = errors.New("save affected row 0")
	UPDATE_AFFECTED_ZERO_ERROR = errors.New("update affected row 0")
	DELETE_AFFECTED_ZERO_ERROR = errors.New("delete affected row 0")
)

func setPage(queryParams map[string]interface{}, db *gorm.DB) *gorm.DB {
	if val, ok := queryParams["pageSize"]; ok {
		pageSize, _ := val.(int)
		if pageSize {{.lt}}= 0 {
			pageSize = 20
		}
        db = db.Limit(pageSize)
		delete(queryParams, "pageSize")

		if val, ok := queryParams["page"]; ok {
			pageNum, _ := val.(int)
			if pageNum {{.lt}}= 1 {
				pageNum = 1
			}
            db = db.Offset((pageNum - 1) * pageSize)
			delete(queryParams, "page")
		}
	}

	return db
}

func getTotal(queryParams map[string]interface{}, db *gorm.DB) int {
	var total int
    db.Where(queryParams).Count(&total)
	return total
}

// test env
func init() {
	dsn := mysql.Config{
		Addr:    config.GetValue("database", "server_address"),
		User:    config.GetValue("database", "user_name"),
		Passwd:  config.GetValue("database", "password"),
		Net:     "tcp",
		DBName:  config.GetValue("database", "db_name"),
		Params:  map[string]string{"charset": "utf8", "parseTime": "True", "loc": "Local"},
		Timeout: time.Duration(5 * time.Second),
	}

	db, err := gorm.Open("mysql", dsn.FormatDSN())
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}

	if err != nil {
		panic(err.Error())
	}

	DbInstance = db
}