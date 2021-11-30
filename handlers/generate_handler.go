package handlers

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"go-generator/log"
	"go-generator/modules"
	"go-generator/util"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type (
	GenerateHandler struct {
		CommonHandler
	}
	GenerateFile struct {
		PackageName  string // 包名
		TemplateName string // 模板名
		FileName     string // 文件名
		Extension    string // 扩展名
	}
)

var (
	GenerateCodeHandler = new(GenerateHandler)
	generateSingleFiles []GenerateFile
)

func init() {
	GenerateCodeHandler.getMapping("code", GenerateCodeHandler.code)

	// 只生成一个文件
	generateSingleFiles = []GenerateFile{
		{PackageName: "common", TemplateName: "common/base", FileName: "base", Extension: "go"},
		{PackageName: "config", TemplateName: "config/config", FileName: "config", Extension: "go"},
		{PackageName: "constant", TemplateName: "constant/error_code", FileName: "error_code", Extension: "go"},
		{PackageName: "control", TemplateName: "control/main_control", FileName: "main_control", Extension: "go"},
		{PackageName: "dto", TemplateName: "dto/common_request", FileName: "common_request", Extension: "go"},
		{PackageName: "dto", TemplateName: "dto/common_result", FileName: "common_result", Extension: "go"},
		{PackageName: "etc", TemplateName: "etc/conf", FileName: "conf", Extension: "ini"},
		{PackageName: "handlers", TemplateName: "handlers/common_handler", FileName: "common_handler", Extension: "go"},
		{PackageName: "log", TemplateName: "log/log", FileName: "log", Extension: "go"},
		{PackageName: "modules", TemplateName: "modules/common_module", FileName: "common_module", Extension: "go"},
		{PackageName: "", TemplateName: "build_bat", FileName: "build", Extension: "bat"},
		{PackageName: "", TemplateName: "build_sh", FileName: "build", Extension: "sh"},
		{PackageName: "", TemplateName: "main", FileName: "main", Extension: "go"},
		{PackageName: "", TemplateName: "gitignore", FileName: "", Extension: "gitignore"},
		{PackageName: "", TemplateName: "readme", FileName: "readme", Extension: "md"},
		{PackageName: "", TemplateName: "api_md", FileName: "接口文档", Extension: "md"},
	}
}

// 生成代码
func (handler *GenerateHandler) code(cont echo.Context) error {
	dbAddr := cont.QueryParam("dbAddr")
	dbUser := cont.QueryParam("dbUser")
	dbPassword := cont.QueryParam("dbPassword")
	dbName := cont.QueryParam("dbName")
	tableNames := cont.QueryParam("tableName")
	projectName := cont.QueryParam("projectName")
	jsonFormat := cont.QueryParam("jsonFormat")
	htmlTemplate := cont.QueryParam("htmlTemplate")

	zipPath := projectName
	if projectName == "" {
		zipPath = dbName
	}
	templateMap := make(map[string]interface{}, 0)
	templateMap["lt"] = template.HTML("<")
	templateMap["packagePrefix"] = zipPath
	templateMap["jsonFormat"] = jsonFormat
	templateMap["handlers"] = []string{}
	// 项目结构 readme文档使用
	projectNames := []string{}
	zipPaths := strings.Split(zipPath, "/")
	spaceStr := ""
	for index, path := range zipPaths {
		if index == 0 {
			projectNames = append(projectNames, spaceStr+path)
		} else {
			projectNames = append(projectNames, spaceStr+"|--"+path)
		}
		spaceStr += "    "
	}
	templateMap["projectNames"] = projectNames

	errorMap := map[string]interface{}{
		"dbAddr":     dbAddr,
		"dbUser":     dbUser,
		"dbPassword": dbPassword,
		"dbName":     dbName,
		"msg":        "连接数据库失败，请检查配置是否正确",
	}
	err := createPath(zipPath, dbAddr)
	if err != nil {
		log.Error("[代码生成][创建文件夹]", "err", err)
		errorMap["msg"] = err.Error()
		return cont.Render(http.StatusOK, htmlTemplate+".html", errorMap)
	}

	if dbAddr != "" {
		db, err := modules.GetDbInstance(dbAddr, dbUser, dbPassword, "")
		defer db.Close()
		if err != nil {
			log.Error("[代码生成][查询数据库]", "err", err)
			return cont.Render(http.StatusOK, htmlTemplate+".html", errorMap)
		}
		templateMap["dbAddr"] = dbAddr
		templateMap["dbUser"] = dbUser
		templateMap["dbPassword"] = dbPassword
		templateMap["dbName"] = dbName
		err = resolveTable(zipPath, db, dbName, tableNames, jsonFormat, templateMap)
		if err != nil {
			log.Error("[代码生成][处理表]", "err", err)
			errorMap["msg"] = err.Error()
			removePath(zipPath, "")
			return cont.Render(http.StatusOK, htmlTemplate+".html", errorMap)
		}
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(generateSingleFiles))
	errorCount := 0
	for _, file := range generateSingleFiles {
		go func(file GenerateFile) {
			defer waitGroup.Done()
			if dbAddr == "" && file.TemplateName == "modules/common_module" {
				return
			}
			err = generateFromTemplate(zipPath, file.PackageName, file.TemplateName, file.FileName, file.Extension, templateMap)
			if err != nil {
				log.Error("[代码生成][生成单个文件失败]", "err", err)
				errorCount++
			}
		}(file)
	}
	waitGroup.Wait()
	if errorCount > 0 {
		errorMap["msg"] = "生成单个文件失败"
		removePath(zipPath, "")
		return cont.Render(http.StatusOK, htmlTemplate+".html", errorMap)
	}

	zipFile := projectNames[0] + ".zip"
	buffer, err := createZip(zipPath, zipFile)
	if err != nil {
		errorMap["msg"] = err.Error()
		removePath(zipPath, "")
		return cont.Render(http.StatusOK, htmlTemplate+".html", errorMap)
	}

	removePath(zipPath, zipFile)
	//设置请求头  使用浏览器下载
	cont.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+zipFile)
	return cont.Stream(http.StatusOK, echo.MIMEOctetStream, bytes.NewReader(buffer))
}

func createPath(path, dbAddr string) error {
	var childPaths []string
	if dbAddr == "" {
		childPaths = []string{"control", "common", "config", "log", "etc", "handlers", "dto", "service", "constant"}
	} else {
		childPaths = []string{"control", "common", "config", "log", "etc", "handlers", "dto", "service", "constant", "modules"}
	}

	for _, childPath := range childPaths {
		err := os.MkdirAll(path+"/"+childPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("创建文件夹失败,path=%s", path+"/"+childPath)
		}
	}
	return nil
}

// 处理表
func resolveTable(zipPath string, db *gorm.DB, dbName, tableNames, jsonFormat string, templateMap map[string]interface{}) error {
	var err error
	tables := []modules.Table{}
	tableArr := []string{}
	if tableNames == "" {
		tables, err = modules.GetTables(db, dbName)
		for _, table := range tables {
			tableArr = append(tableArr, table.TableName)
		}
		if err != nil {
			log.Error("[代码生成][查询数据库表]", "err", err)
		}
	} else {
		tableArr = strings.Split(tableNames, ",")
		tables, err = modules.GetTablesIn(db, dbName, tableArr)
		if err != nil {
			log.Error("[代码生成][查询数据库表]", "err", err)
		}
	}
	// 查询所有表的字段 然后组装
	columns, err := modules.GetColumnsIn(db, dbName, tableArr)
	if err != nil {
		log.Error("[代码生成][查询表字段失败]", "err", err)
		return fmt.Errorf("查询表字段失败, error=%s", err.Error())
	}
	for index, table := range tables {
		for _, column := range columns {
			if column.TableName == table.TableName {
				tables[index].Columns = append(tables[index].Columns, column)
			}
		}
	}
	handlers := []string{}
	tableColumns := []modules.TableColumn{}
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(tables))
	errorCount := 0
	for index, table := range tables {
		go func(index int, table modules.Table) {
			structName := util.ToUpperCamelCase(table.TableName)
			tableComment := ""
			if table.TableComment == "" {
				tableComment = table.TableName
			} else {
				tableComment = strings.Replace(table.TableComment, "表", "", -1)
			}
			urlPrefix := strings.Replace(table.TableName, "_", "-", -1)
			tableColumn := modules.TableColumn{
				Columns:          columns,
				NoPkColumns:      columns[1:],
				Struct:           structName,
				LowerCamelStruct: util.ToLowerCamelCase(structName),
				UrlPrefix:        urlPrefix,
				TableName:        table.TableName,
				TableComment:     tableComment,
				PackagePrefix:    zipPath,
				Index:            index + 1,
				JsonFormat:       jsonFormat,
			}

			// 处理表字段
			tableColumn = resolveColumn(table.Columns, jsonFormat, tableColumn)
			tables[index].TableColumn = tableColumn

			columnMap := make(map[string]modules.Column, len(table.Columns))
			for _, column := range tableColumn.Columns {
				columnMap[column.ColumnName] = column
			}
			// 查询单表索引，为生成新增/修改接口重复校验代码做准备
			tableIndexes, err := modules.GetIndex(db, dbName, table.TableName)
			if err != nil {
				errorCount++
				log.Error("[代码生成][查询单表索引失败]", "err", err)
			}
			tableIndexMap := make(map[string][]modules.Column, len(tableIndexes))
			for _, tableIndex := range tableIndexes {
				if tableIndex.NonUnique == 1 {
					continue
				} else if tableIndex.KeyName == "PRIMARY" {
					continue
				}
				if _, ok := tableIndexMap[tableIndex.KeyName]; ok {
					tableIndexMap[tableIndex.KeyName] = append(tableIndexMap[tableIndex.KeyName], columnMap[tableIndex.ColumnName])
				} else {
					tableIndexMap[tableIndex.KeyName] = []modules.Column{columnMap[tableIndex.ColumnName]}
				}
			}
			tableColumn.HasUniqueIndex = false
			if len(tableIndexMap) > 0 {
				indexes := make([]modules.Index, 0, len(tableIndexMap))
				for key, columns := range tableIndexMap {
					indexes = append(indexes, modules.Index{
						KeyName: key,
						Columns:columns,
					})
				}
				tableColumn.Indexes = indexes
				tableColumn.HasUniqueIndex = true
			}

			// 从模板文件生成dto代码文件
			err = generateFromTemplate(zipPath, "dto", "dto/dto", table.TableName, "go", tableColumn)
			if err != nil {
				errorCount++
				log.Error("[代码生成][生成dto文件失败]", "err", err)
				//return fmt.Errorf("生成handler文件失败, error=%s", err.Error())
			}

			// 从模板文件生成handler代码文件
			err = generateFromTemplate(zipPath, "handlers", "handlers/handlers", table.TableName+"_handler", "go", tableColumn)
			if err != nil {
				errorCount++
				log.Error("[代码生成][生成handler文件失败]", "err", err)
				//return fmt.Errorf("生成handler文件失败, error=%s", err.Error())
			}

			// 从模板文件生成module代码文件
			err = generateFromTemplate(zipPath, "modules", "modules/modules", table.TableName+"_module", "go", tableColumn)
			if err != nil {
				errorCount++
				log.Error("[代码生成][生成module文件失败]", "err", err)
				//return fmt.Errorf("生成module文件失败, error=%s", err.Error())
			}

			// 从模板文件生成handler代码文件
			err = generateFromTemplate(zipPath, "service", "service/service", table.TableName+"_service", "go", tableColumn)
			if err != nil {
				errorCount++
				log.Error("[代码生成][生成service文件失败]", "err", err)
				//return fmt.Errorf("生成handler文件失败, error=%s", err.Error())
			}

			fmtCode(zipPath + "/dto" + "/" + table.TableName + ".go")
			fmtCode(zipPath + "/handlers" + "/" + table.TableName + "_handler.go")
			fmtCode(zipPath + "/modules" + "/" + table.TableName + "_module.go")
			fmtCode(zipPath + "/service" + "/" + table.TableName + "_service.go")
			defer waitGroup.Done()
		}(index, table)
	}
	waitGroup.Wait()
	if errorCount > 0 {
		return fmt.Errorf("处理表失败")
	}
	for _, table := range tables {
		handlers = append(handlers, "Default"+table.TableColumn.Struct+"Handler", "Page"+table.TableColumn.Struct+"Handler")
		tableColumns = append(tableColumns, table.TableColumn)
	}
	templateMap["tableColumns"] = tableColumns
	templateMap["handlers"] = handlers
	return nil
}

// 处理表字段
func resolveColumn(columns []modules.Column, jsonFormat string, tableColumn modules.TableColumn) modules.TableColumn {
	tableColumn.HasTime = false
	for i, column := range columns {
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

		upperColumn := util.ToUpperCamelCase(column.ColumnName)
		columns[i].UpperColumn = upperColumn
		lowerColumn := util.ToLowerCamelCase(upperColumn)
		columns[i].LowerColumn = lowerColumn
		if jsonFormat == "1" {
			columns[i].JsonColumn = column.ColumnName
		} else {
			columns[i].JsonColumn = lowerColumn
		}

		if i == 0 {
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
		}
	}
	tableColumn.Columns = columns
	tableColumn.NoPkColumns = columns[1:]
	return tableColumn
}

// 从模板生成
func generateFromTemplate(zipPath, packageName, templateName, fileName, extension string, templateParam interface{}) error {
	var t *template.Template
	t, err := template.ParseFiles("templates/generate/" + templateName + ".tpl")
	if t == nil || err != nil {
		return fmt.Errorf("读取模板文件失败,error=%s", err.Error())
	}
	packagePath := zipPath
	if packageName != "" {
		packagePath += "/" + packageName
	}

	file := packagePath + "/" + fileName + "." + extension
	exists, _ := util.PathExists(file)
	var f *os.File
	if exists {
		f, err = os.Open(file)
		if err != nil {
			return fmt.Errorf("打开文件失败,file=%s", file)
		}
	} else {
		f, err = os.Create(file)
		if err != nil {
			return fmt.Errorf("创建文件失败,file=%s", file)
		}
	}

	err = t.Execute(f, templateParam)
	f.Close()
	if err != nil {
		log.Error("[代码生成][从模板文件生成失败]", "err", err)
		return fmt.Errorf("从模板文件生成失败,error=%s", err.Error())
	}
	return nil
}

// gofmt 格式化go代码
func fmtCode(file string) {
	cmd := exec.Command("go", "fmt", file)
	if err := cmd.Run(); err != nil {
		log.Error("[代码生成][格式化代码失败]", "err", err, "file", file)
		//return fmt.Errorf("格式化代码失败,error=%s,file=%s", err.Error(), file)
	}
}

// 创建zip
func createZip(zipPath, zipFile string) ([]byte, error) {
	zipPaths := strings.Split(zipPath, "/")
	err := util.Zip(zipFile, zipPaths[0])
	if err != nil {
		return nil, fmt.Errorf("压缩zip包失败,error=%s", err.Error())
	}
	zip, err := os.Open(zipFile)
	if err != nil {
		return nil, err
	}

	fileInfo, err := zip.Stat()
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, fileInfo.Size())
	_, err = zip.Read(buffer)
	if err != nil {
		return nil, err
	}
	defer zip.Close()
	return buffer, nil
}

// 删除临时文件夹和压缩包
func removePath(zipPath, zipFile string) error {
	zipPaths := strings.Split(zipPath, "/")
	err := os.RemoveAll(zipPaths[0])
	if err != nil {
		return fmt.Errorf("删除文件夹失败,error=%s", err.Error())
	}
	if zipFile != "" {
		err = os.Remove(zipFile)
		if err != nil {
			return fmt.Errorf("删除zip包失败,error=%s", err.Error())
		}
	}
	return nil
}
