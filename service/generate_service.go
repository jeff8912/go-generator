package service

import (
	"go-generator/constant"
	"go-generator/dto"
	"go-generator/log"
	"go-generator/module"
	"go-generator/util"
	"gorm.io/gorm"
	"html/template"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type (
	GenerateServiceImpl struct {
		CommonServiceImpl
	}
)

var (
	GenerateService  = new(GenerateServiceImpl)
	templateRootPath = "template/generate/"
	singleFiles      []string
	multiFiles       []string
)

func init() {
	initTplFiles()
}

func (p *GenerateServiceImpl) Generate(param *dto.GenerateParam) error {
	templateMap := make(map[string]interface{}, 0)
	templateMap["lt"] = template.HTML("<")
	templateMap["packagePrefix"] = param.ProjectName
	templateMap["Now"] = template.HTML(`{{ .Now "2006-01-02 15:04:05" }}`)

	if param.DataSource == constant.DATASOURCE_MYSQL {
		templateMap["dbAddr"] = param.DbAddr
		templateMap["dbUser"] = param.Username
		templateMap["dbPassword"] = param.Password
		templateMap["dbName"] = param.DbName

		// 获取数据库连接
		db, err := module.GetDbInstance(param.DbAddr, param.Username, param.Password, "")
		if err != nil {
			return err
		}

		// 获取表列表和表名列表
		tables, tableNames, err := getTables(db, param)
		if err != nil {
			return err
		}

		// 初始化表及其列map
		tableColumnMap := initTableColumnMap(tables, param)

		// 加载表字段
		columnMap, err := loadColumns(db, tableColumnMap, tableNames, param)
		if err != nil {
			return err
		}

		// 加载表索引
		err = loadIndexes(db, tableColumnMap, columnMap, param)
		if err != nil {
			return err
		}

		// 按表生成多个文件
		err = generateMultiFiles(tableColumnMap, param)
		if err != nil {
			return err
		}

		// 接口文档生成需要的数据
		tableColumns := make([]*dto.TableColumn, 0, len(tables))
		for i := range tables {
			tableColumns = append(tableColumns, tableColumnMap[tables[i].TableName])
		}
		templateMap["tableColumns"] = tableColumns
	}

	err := generateSingleFiles(param, templateMap)
	if err != nil {
		return err
	}

	return nil
}

func getTables(db *gorm.DB, param *dto.GenerateParam) (tables []*dto.TableOptions, tableNames []string, err error) {
	if len(param.TableNames) <= 0 {
		tables, err = module.GetTables(db, param.DbName)
		for _, table := range tables {
			tableNames = append(tableNames, table.TableName)
		}
		if err != nil {
			log.Error("[代码生成][查询数据库表]", "err", err)
			return nil, nil, err
		}
	} else {
		tableNames = param.TableNames
		tables, err = module.GetTablesIn(db, param.DbName, tableNames)
		if err != nil {
			log.Error("[代码生成][查询数据库表]", "err", err)
			return nil, nil, err
		}
	}
	return tables, tableNames, nil
}

func initTableColumnMap(tables []*dto.TableOptions, param *dto.GenerateParam) map[string]*dto.TableColumn {
	tableColumnMap := make(map[string]*dto.TableColumn, len(tables))
	for index, table := range tables {
		structName := util.ToUpperCamelCase(table.TableName)
		tableComment := ""
		if table.TableComment == "" {
			tableComment = table.TableName
		} else {
			tableComment = strings.Replace(table.TableComment, "表", "", -1)
		}
		urlPrefix := strings.Replace(table.TableName, "_", "-", -1)
		tableColumn := &dto.TableColumn{
			HasTime:          false,
			Struct:           structName,
			LowerCamelStruct: util.ToLowerCamelCase(structName),
			UrlPrefix:        urlPrefix,
			TableName:        table.TableName,
			TableComment:     tableComment,
			PackagePrefix:    param.ProjectName,
			Index:            index + 1,
		}
		tableColumnMap[table.TableName] = tableColumn
	}
	return tableColumnMap
}

func loadColumns(db *gorm.DB, tableColumnMap map[string]*dto.TableColumn, tableNames []string, param *dto.GenerateParam) (
	columnMap map[string]*dto.Column, err error) {
	// 查询所有表的字段 然后组装
	columns, err := module.GetColumnsIn(db, param.DbName, tableNames)
	if err != nil {
		log.Error("[代码生成][查询表字段失败]", "err", err)
		return columnMap, err
	}

	columnMap = make(map[string]*dto.Column, len(columns))
	for i, column := range columns {
		tableColumn := tableColumnMap[column.TableName]
		columnMap[column.ColumnName] = columns[i]

		columns[i].ColumnName = strings.ToLower(column.ColumnName)
		columns[i].DataType = strings.ToLower(column.DataType)
		// 是否包含时间
		if column.DataType == "time" || column.DataType == "date" || column.DataType == "year" ||
			column.DataType == "timestamp" || column.DataType == "datetime" {
			tableColumn.HasTime = true
			columns[i].ColumnType = "time.Time"

		} else if column.DataType == "char" || column.DataType == "varchar" || column.DataType == "tinyblob" ||
			column.DataType == "tinytext" || column.DataType == "blob" || column.DataType == "text" ||
			column.DataType == "mediumblob" || column.DataType == "mediumtext" ||
			column.DataType == "longblob" || column.DataType == "longtext" {
			columns[i].ColumnType = "string"

		} else if column.DataType == "tinyint" || column.DataType == "smallint" || column.DataType == "mediumint" ||
			column.DataType == "int" || column.DataType == "integer" {
			columns[i].ColumnType = "int"

		} else if column.DataType == "bigint" {
			columns[i].ColumnType = "int64"

		} else if column.DataType == "float" {
			columns[i].ColumnType = "float32"

		} else if column.DataType == "double" || column.DataType == "decimal" || column.DataType == "numeric" {
			columns[i].ColumnType = "float64"
		}

		// 获取主键
		upperColumn := util.ToUpperCamelCase(column.ColumnName)
		columns[i].UpperColumn = upperColumn
		lowerColumn := util.ToLowerCamelCase(upperColumn)
		columns[i].LowerColumn = lowerColumn
		columns[i].JsonColumn = column.ColumnName

		if len(tableColumn.Columns) == 0 {
			tableColumn.Pk = column.ColumnName
			tableColumn.LowerCamelPk = columns[i].LowerColumn
			tableColumn.UpperCamelPk = columns[i].UpperColumn
			if column.DataType == "varchar" {
				tableColumn.PkType = "string"
			} else if column.DataType == "bigint" {
				tableColumn.PkType = "int64"
			} else {
				tableColumn.PkType = "int"
			}
		} else {
			tableColumn.NoPkColumns = append(tableColumn.NoPkColumns, columns[i])
		}
		tableColumn.Columns = append(tableColumn.Columns, columns[i])
	}
	return columnMap, err
}

func loadIndexes(db *gorm.DB, tableColumnMap map[string]*dto.TableColumn, columnMap map[string]*dto.Column,
	param *dto.GenerateParam) (err error) {
	for tableName := range tableColumnMap {
		tableColumn := tableColumnMap[tableName]
		// 查询单表索引，为生成新增/修改接口重复校验代码做准备
		tableIndexes, err := module.GetIndex(db, param.DbName, tableName)
		if err != nil {
			return err
		}
		tableIndexMap := make(map[string][]dto.Column, len(tableIndexes))
		for _, tableIndex := range tableIndexes {
			if tableIndex.NonUnique == 1 {
				continue
			} else if tableIndex.KeyName == "PRIMARY" {
				continue
			}
			tableIndexMap[tableIndex.KeyName] = append(tableIndexMap[tableIndex.KeyName], *columnMap[tableIndex.ColumnName])
		}
		tableColumn.HasUniqueIndex = false
		if len(tableIndexMap) > 0 {
			indexes := make([]dto.Index, 0, len(tableIndexMap))
			for key, columns := range tableIndexMap {
				indexes = append(indexes, dto.Index{
					KeyName: key,
					Columns: columns,
				})
			}
			tableColumn.Indexes = indexes
			tableColumn.HasUniqueIndex = true
		}
	}
	return nil
}

func generateMultiFiles(tableColumnMap map[string]*dto.TableColumn, param *dto.GenerateParam) error {
	var waitGroup sync.WaitGroup
	waitSize := len(tableColumnMap) * len(multiFiles)
	waitGroup.Add(waitSize)
	errs := make([]error, 0, waitSize)
	for tableName := range tableColumnMap {
		for _, file := range multiFiles {
			go func(tableName, file string) {
				defer waitGroup.Done()
				err := generateFromTemplate(param.ProjectName, file, tableName, tableColumnMap[tableName])
				if err != nil {
					log.Error("[代码生成][生成多个文件失败]", "err", err)
				}
				fmtCode(param.ProjectName + "/" + file)
			}(tableName, file)
		}
	}
	waitGroup.Wait()
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func generateSingleFiles(param *dto.GenerateParam, templateMap map[string]interface{}) error {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(singleFiles))
	errs := make([]error, 0, len(singleFiles))
	for _, file := range singleFiles {
		go func(file string) {
			defer waitGroup.Done()
			if param.DataSource == constant.DATASOURCE_NONE && strings.Contains(file, "module/") {
				return
			}
			err := generateFromTemplate(param.ProjectName, file, "", templateMap)
			if err != nil {
				log.Error("[代码生成][生成单个文件失败]", "err", err)
				errs = append(errs, err)
			}
		}(file)
	}
	waitGroup.Wait()
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

// 从模板生成
func generateFromTemplate(zipPath, templateName, tableName string, templateParam interface{}) error {
	var t *template.Template
	t, err := template.ParseFiles(templateRootPath + templateName)
	if err != nil {
		return err
	}

	file := strings.ReplaceAll(templateName, ".tpl", "")
	if tableName != "" {
		file = strings.ReplaceAll(file, "/", "/"+tableName+"_")
	}
	file = zipPath + "/" + file

	f, err := util.CreateFileIfAbsent(file)
	if err != nil {
		return err
	}

	err = t.Execute(f, templateParam)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

// gofmt 格式化go代码
func fmtCode(file string) {
	file = strings.ReplaceAll(file, ".tpl", "")
	cmd := exec.Command("go", "fmt", file)
	if err := cmd.Run(); err != nil {
		log.Error("[代码生成][格式化代码失败]", "err", err, "file", file)
	}
}

func initTplFiles() {
	files, err := util.GetAllFile(templateRootPath)
	if err != nil {
		panic(err)
	}

	for i := range files {
		file := strings.ReplaceAll(files[i], string(filepath.Separator), "/")
		file = strings.ReplaceAll(file, templateRootPath, "")
		if strings.Contains(file, ".tpl") {
			if strings.Contains(file, "/dto.go.tpl") ||
				strings.Contains(file, "/module.go.tpl") ||
				strings.Contains(file, "/service.go.tpl") ||
				strings.Contains(file, "/handler.go.tpl") {
				multiFiles = append(multiFiles, file)
			} else {
				singleFiles = append(singleFiles, file)
			}
		}
	}
}
