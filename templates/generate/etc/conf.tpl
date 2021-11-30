[system]
server_address = 127.0.0.1:8080
local_ip = 127.0.0.1
hostname = 127.0.0.1

{{if .dbAddr}}[database]
user_name = {{.dbUser}}
password = {{.dbPassword}}
server_address = {{.dbAddr}}
db_name = {{.dbName}}{{end}}

[log]
level = debug

