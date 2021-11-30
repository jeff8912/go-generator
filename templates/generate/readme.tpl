## 编译与打包脚本说明：
	build.bat：windows下编译脚本，注：暂时不支持打包命令。
	build.sh：linux编译脚本，支持打包，打包命令：build.sh pkg

## 工程说明：
```
目录说明：
	1.control目录：存放工程启动调度go文件
	2.etc目录：存放工程配置文件
	3.handlers目录：存放业务处理go文件
	4.modules目录：存放业务数据访问go文件
	5.log目录：存放系统运行日志

文件说明：
	1.main.go文件：工程主入口文件
	2.main_control.go文件：工程主调度文件
	3.build.sh文件：linux下工程编译脚本
	4.build.bat文件：window下工程编译脚本
	5.conf.ini文件：工程默认配置文件
```

## 项目结构:
```
内部公共组名：          os_comm
    |--内部公共工程名：     |--os_go_comm
    |--内部公共工程名：  	 |--vendor

私有项目结构：          {{ range $i, $v := .projectNames }}{{$v}}
					 {{end}}
项目开发目录：          src
                         |--os_go_comm
                         |--vendor
						 {{ range $i, $v := .projectNames }}{{if eq $i 0}}|--{{end}}{{$v}}
					 	 {{end}}
```
