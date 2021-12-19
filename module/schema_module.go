package module

import (
	"fmt"
	"go-generator/dto"
	"gorm.io/gorm"
)

func GetDatabases(dbInstance *gorm.DB) ([]dto.Database, error) {
	records := []dto.Database{}

	var sql string
	sql = `SELECT DISTINCT table_schema AS table_schema FROM information_schema.tables ORDER BY table_schema`

	res := dbInstance.Debug().Raw(sql).Scan(&records)
	if res.Error != nil {
		return records, res.Error
	}

	return records, nil
}

func GetTables(dbInstance *gorm.DB, dbName string) ([]*dto.TableOptions, error) {
	records := []*dto.TableOptions{}

	var sql string
	sql = `SELECT table_name AS table_name, table_comment AS table_comment FROM information_schema.tables WHERE table_schema = ?`

	res := dbInstance.Debug().Raw(sql, dbName).Scan(&records)
	if res.Error != nil {
		return records, res.Error
	}

	return records, nil
}

func GetTablesIn(dbInstance *gorm.DB, dbName string, tableArr []string) ([]*dto.TableOptions, error) {
	records := []*dto.TableOptions{}

	var sql string
	sql = `SELECT table_name AS table_name, table_comment AS table_comment FROM information_schema.tables WHERE table_schema = ? AND table_name IN (?)`

	res := dbInstance.Debug().Raw(sql, dbName, tableArr).Scan(&records)
	if res.Error != nil {
		return records, res.Error
	}

	return records, nil
}

func GetColumns(dbInstance *gorm.DB, dbName, tableName string) ([]*dto.Column, error) {
	records := []*dto.Column{}

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

func GetColumnsIn(dbInstance *gorm.DB, dbName string, tableArr []string) ([]*dto.Column, error) {
	records := []*dto.Column{}

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

func GetIndex(dbInstance *gorm.DB, dbName string, tableName string) ([]*dto.TableIndex, error) {
	records := []*dto.TableIndex{}

	res := dbInstance.Debug().Raw(fmt.Sprintf("SHOW INDEX FROM %s.%s", dbName, tableName)).Scan(&records)
	if res.Error != nil {
		return records, res.Error
	}

	return records, nil
}
