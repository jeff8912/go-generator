package handlers

import (
	"github.com/labstack/echo"
	"go-generator/config"
	"go-generator/log"
	"go-generator/modules"
	"net/http"
)

var (
	DefaultIndexHandler            = new(TemplateHandler)
	DefaultConfigHandler           = new(TemplateHandler)
	DefaultGenerate4MysqlHandler   = new(TemplateHandler)
	DefaultGenerate4NoMysqlHandler = new(TemplateHandler)
	DefaultDatabaseHandler         = new(TemplateHandler)
	DefaultTableHandler            = new(TemplateHandler)
)

type (
	TemplateHandler struct {
		CommonHandler
	}
)

var (
	Addr     string
	User     string
	Password string
	DBName   string
)

func init() {
	Addr = config.GetValue("database", "server_address")
	User = config.GetValue("database", "user_name")
	Password = config.GetValue("database", "password")
	DBName = config.GetValue("database", "db_name")

	DefaultIndexHandler.getMapping("", DefaultIndexHandler.config)
	DefaultConfigHandler.getMapping("config", DefaultConfigHandler.config)
	DefaultGenerate4MysqlHandler.getMapping("generate4mysql", DefaultGenerate4MysqlHandler.generate4mysql)
	DefaultGenerate4NoMysqlHandler.getMapping("generate4nomysql", DefaultGenerate4NoMysqlHandler.generate4nomysql)
	DefaultDatabaseHandler.getMapping("database", DefaultDatabaseHandler.database)
	DefaultTableHandler.getMapping("table", DefaultTableHandler.table)
}

func (handler *TemplateHandler) config(cont echo.Context) error {
	return cont.Render(http.StatusOK, "config.html", map[string]interface{}{
		"dbAddr":     Addr,
		"dbUser":     User,
		"dbPassword": Password,
		"dbName":     DBName,
	})
}

func (handler *TemplateHandler) generate4mysql(cont echo.Context) error {
	dbAddr := cont.QueryParam("dbAddr")
	dbUser := cont.QueryParam("dbUser")
	dbPassword := cont.QueryParam("dbPassword")
	//dbName := cont.QueryParam("dbName")

	db, err := modules.GetDbInstance(dbAddr, dbUser, dbPassword, "")

	if err != nil {
		log.Error("[确认配置][检查数据库配置是否正确]", "err", err)
		return cont.Render(http.StatusOK, "config.html", map[string]interface{}{
			"dbAddr":     dbAddr,
			"dbUser":     dbUser,
			"dbPassword": dbPassword,
			"dbName":     DBName,
			"msg":        "连接数据库失败，请检查配置是否正确",
		})
	}
	defer db.Close()

	return cont.Render(http.StatusOK, "generate4mysql.html", map[string]interface{}{
		"dbAddr":     dbAddr,
		"dbUser":     dbUser,
		"dbPassword": dbPassword,
		"dbName":     DBName,
	})
}

func (handler *TemplateHandler) generate4nomysql(cont echo.Context) error {
	return cont.Render(http.StatusOK, "generate4nomysql.html", nil)
}

func (handler *TemplateHandler) database(cont echo.Context) error {
	dbAddr := cont.QueryParam("dbAddr")
	dbUser := cont.QueryParam("dbUser")
	dbPassword := cont.QueryParam("dbPassword")
	//dbName := cont.QueryParam("dbName")

	db, err := modules.GetDbInstance(dbAddr, dbUser, dbPassword, "")

	if err != nil {
		log.Error("[代码生成][查询数据库]", "err", err)
		return cont.Render(http.StatusOK, "config.html", map[string]interface{}{
			"dbAddr":     dbAddr,
			"dbUser":     dbUser,
			"dbPassword": dbPassword,
			"dbName":     DBName,
			"msg":        "连接数据库失败，请检查配置是否正确",
		})
	}
	databases, err := modules.GetDatabases(db)
	if err != nil {
		log.Error("[代码生成][查询数据库]", "err", err)
		return cont.Render(http.StatusOK, "config.html", map[string]interface{}{
			"dbAddr":     dbAddr,
			"dbUser":     dbUser,
			"dbPassword": dbPassword,
			"dbName":     DBName,
			"msg":        "查询数据库失败，请检查配置是否正确",
		})
	}
	defer db.Close()
	res := BaseResEntity{}
	res.StatusDesc = "操作成功!"
	res.Data = databases
	return cont.JSON(http.StatusOK, res)
}

func (handler *TemplateHandler) table(cont echo.Context) error {
	dbAddr := cont.QueryParam("dbAddr")
	dbUser := cont.QueryParam("dbUser")
	dbPassword := cont.QueryParam("dbPassword")
	dbName := cont.QueryParam("dbName")

	db, err := modules.GetDbInstance(dbAddr, dbUser, dbPassword, dbName)

	if err != nil {
		log.Error("[代码生成][查询数据库表]", "err", err)
		return cont.Render(http.StatusOK, "config.html", map[string]interface{}{
			"dbAddr":     dbAddr,
			"dbUser":     dbUser,
			"dbPassword": dbPassword,
			"dbName":     dbName,
			"msg":        "连接数据库失败，请检查配置是否正确",
		})
	}
	tables, err := modules.GetTables(db, dbName)
	if err != nil {
		log.Error("[代码生成][查询数据库表]", "err", err)
		return cont.Render(http.StatusOK, "config.html", map[string]interface{}{
			"dbAddr":     dbAddr,
			"dbUser":     dbUser,
			"dbPassword": dbPassword,
			"dbName":     dbName,
			"msg":        "查询数据库表失败，请检查配置是否正确",
		})
	}
	defer db.Close()
	res := BaseResEntity{}
	res.StatusDesc = "操作成功!"
	res.Data = tables
	return cont.JSON(http.StatusOK, res)
}
