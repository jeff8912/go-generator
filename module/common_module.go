package module

import (
	"fmt"
	"go-generator/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var DbInstance *gorm.DB

func GetDbInstance(dbAddr, dbUser, dbPassword, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbUser, dbPassword, dbAddr, dbName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(config.Int("database", "max_idle_conns"))
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(config.Int("database", "max_open_conns"))
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Duration(config.Int("database", "conn_max_lifetime")) * time.Hour)

	DbInstance = db
	return DbInstance, nil
}
