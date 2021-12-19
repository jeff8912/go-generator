## 编译与打包脚本说明：
	build.bat：windows下编译脚本，注：暂时不支持打包命令。
	build.sh：linux编译脚本，支持打包，打包命令：build.sh pkg

## 工程说明：
```
目录说明：
    1. banner：banner以及logo打印
    2. config：配置文件读取
    3. constant：常量
    4. control：工程启动调度
    5. dto：数据传输对象
    6. etc：工程配置文件
    7. handler：控制层
    8. log：日志功能
    9. module：数据访问层
    10. service：业务层

文件说明：
	1. main.go文件：工程主入口文件
	2. main_control.go文件：工程主调度文件
	3. build.sh文件：linux下工程编译脚本
	4. build.bat文件：window下工程编译脚本
	5. conf.ini文件：工程默认配置文件
```