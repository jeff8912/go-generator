package modules

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"time"
)

var DbInstance *gorm.DB

var (
	SAVE_AFFECTED_ZERO_ERROR   = errors.New("save affected row 0")
	UPDATE_AFFECTED_ZERO_ERROR = errors.New("update affected row 0")
	DELETE_AFFECTED_ZERO_ERROR = errors.New("delete affected row 0")
)

// test env
func GetDbInstance(dbAddr, dbUser, dbPassword, dbName string) (*gorm.DB, error) {
	dsn := mysql.Config{
		Addr:    dbAddr,
		User:    dbUser,
		Passwd:  dbPassword,
		Net:     "tcp",
		Params:  map[string]string{"charset": "utf8", "parseTime": "True", "loc": "Local"},
		Timeout: time.Duration(5 * time.Second),
	}
	if dbName != "" {
		dsn.DBName = dbName
	}

	db, err := gorm.Open("mysql", dsn.FormatDSN())
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return defaultTableName
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
