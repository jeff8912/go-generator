package modules

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type (
	Database struct {
		Id          int    `json:"id"`
		TableSchema string `json:"text" gorm:"column:table_schema"` // 数据库名称
	}

	Table struct {
		TableName    string `json:"table_name" gorm:"column:table_name"`       // 表名
		TableComment string `json:"table_comment" gorm:"column:table_comment"` // 表注释
		Columns      []Column
		TableColumn  TableColumn
	}

	Column struct {
		TableName     string `json:"table_name" gorm:"column:table_name"`         // 表名
		ColumnName    string `json:"column_name" gorm:"column:column_name"`       // 字段名
		DataType      string `json:"data_type" gorm:"column:data_type"`           // 字段类型
		ColumnComment string `json:"column_comment" gorm:"column:column_comment"` // 字段注释
		IsNullable    string `json:"is_nullable" gorm:"column:is_nullable"`       // 列的为空性。如果列允许 NULL，那么该列返回 YES。否则，返回 NO。
		UpperColumn   string                                                      // 大驼峰字段
		LowerColumn   string                                                      // 小驼峰字段
		JsonColumn    string                                                      // json字段
		ColumnType    string                                                      // 字段类型
	}

	TableColumn struct {
		Columns          []Column `json:"columns"`          // 表字段
		NoPkColumns      []Column `json:"noPkColumns"`      // 不带主键的表字段
		Struct           string   `json:"struct"`           // 表结构体定义
		LowerCamelStruct string   `json:"lowerCamelStruct"` // 驼峰结构体
		UrlPrefix        string   `json:"urlPrefix"`        // restful路径前缀
		TableName        string   `json:"tableName"`        // 表名
		TableComment     string   `json:"tableComment"`     // 表注释
		Index            int      `json:"index"`            // 序号
		Pk               string   `json:"pk"`               // 主键
		PkType           string   `json:"pkType"`           // 主键类型
		LowerCamelPk     string   `json:"lowerCamelPk"`     // 驼峰主键
		UpperCamelPk     string   `json:"upperCamelPk"`     // 大驼峰主键
		HasTime          bool     `json:"hasTime"`          // 是否包含时间字段
		PackagePrefix    string   `json:"packagePrefix"`    // 包前缀
		JsonFormat       string   `json:"jsonFormat"`       // json字段格式 1：下划线 2：驼峰
		HasUniqueIndex   bool     `json:"hasUniqueIndex"`   // 是否包含唯一索引
		Indexes          []Index
	}

	TableIndex struct {
		Table      string `json:"table" gorm:"column:Table"`               // 表名
		NonUnique  int    `json:"non_unique" gorm:"column:Non_unique"`     // 非唯一键
		KeyName    string `json:"key_name" gorm:"column:Key_name"`         // 索引名称
		SeqInIndex string `json:"seq_in_index" gorm:"column:Seq_in_index"` // 在当前索引中的序号
		ColumnName string `json:"column_name" gorm:"column:Column_name"`   // 字段名称
	}

	Index struct {
		KeyName string // 索引名称
		Columns []Column
	}
)

func GetDatabases(dbInstance *gorm.DB) ([]Database, error) {
	records := []Database{}

	var sql string
	sql = `SELECT DISTINCT table_schema AS table_schema FROM information_schema.tables ORDER BY table_schema`

	res := dbInstance.Debug().Raw(sql).Scan(&records)
	if res.Error != nil {
		return records, res.Error
	}

	return records, nil
}

func GetTables(dbInstance *gorm.DB, dbName string) ([]Table, error) {
	records := []Table{}

	var sql string
	sql = `SELECT table_name AS table_name, table_comment AS table_comment FROM information_schema.tables WHERE table_schema = ?`

	res := dbInstance.Debug().Raw(sql, dbName).Scan(&records)
	if res.Error != nil {
		return records, res.Error
	}

	return records, nil
}

func GetTablesIn(dbInstance *gorm.DB, dbName string, tableArr []string) ([]Table, error) {
	records := []Table{}

	var sql string
	sql = `SELECT table_name AS table_name, table_comment AS table_comment FROM information_schema.tables WHERE table_schema = ? AND table_name IN (?)`

	res := dbInstance.Debug().Raw(sql, dbName, tableArr).Scan(&records)
	if res.Error != nil {
		return records, res.Error
	}

	return records, nil
}

func GetColumns(dbInstance *gorm.DB, dbName, tableName string) ([]Column, error) {
	records := []Column{}

	var sql string
	sql = `SELECT
		table_name AS table_name,
		column_name AS column_name,
		data_type AS data_type,
		column_comment AS column_comment,
		is_nullable AS is_nullable
	FROM information_schema.COLUMNS
	WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?`

	res := dbInstance.Debug().Raw(sql, dbName, tableName).Scan(&records)
	if res.Error != nil {
		return records, res.Error
	}

	return records, nil
}

func GetColumnsIn(dbInstance *gorm.DB, dbName string, tableArr []string) ([]Column, error) {
	records := []Column{}

	var sql string
	sql = `SELECT
		table_name AS table_name,
		column_name AS column_name,
		data_type AS data_type,
		column_comment AS column_comment,
		is_nullable AS is_nullable
	FROM information_schema.COLUMNS
	WHERE TABLE_SCHEMA = ? AND TABLE_NAME IN (?)`

	res := dbInstance.Debug().Raw(sql, dbName, tableArr).Scan(&records)
	if res.Error != nil {
		return records, res.Error
	}

	return records, nil
}

func GetIndex(dbInstance *gorm.DB, dbName string, tableName string) ([]TableIndex, error) {
	records := []TableIndex{}

	res := dbInstance.Debug().Raw(fmt.Sprintf("SHOW INDEX FROM %s.%s", dbName, tableName)).Scan(&records)
	if res.Error != nil {
		return records, res.Error
	}

	return records, nil
}
