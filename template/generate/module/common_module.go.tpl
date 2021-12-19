package module

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func getTotal(queryParams map[string]interface{}, db *gorm.DB) int64 {
	var total int64
    db.Where(queryParams).Count(&total)
	return total
}

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
        config.GetValue("database", "user_name"),
        config.GetValue("database", "password"),
        config.GetValue("database", "server_address"),
        config.GetValue("database", "db_name"),
    )
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        PrepareStmt: true,
    })

    sqlDB, err := db.DB()
    if err != nil {
        panic(err.Error())
    }

    // SetMaxIdleConns 设置空闲连接池中连接的最大数量
    sqlDB.SetMaxIdleConns(config.Int("database", "max_idle_conns"))
    // SetMaxOpenConns 设置打开数据库连接的最大数量。
    sqlDB.SetMaxOpenConns(config.Int("database", "max_open_conns"))
    // SetConnMaxLifetime 设置了连接可复用的最大时间。
    sqlDB.SetConnMaxLifetime(time.Duration(config.Int("database", "conn_max_lifetime")) * time.Hour)

    DbInstance = db
}