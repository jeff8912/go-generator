[system]
server_address = 127.0.0.1:8080

{{if .dbAddr}}[database]
user_name = {{.dbUser}}
password = {{.dbPassword}}
server_address = {{.dbAddr}}
db_name = {{.dbName}}
max_idle_conns = 10
max_open_conns = 100
conn_max_lifetime = 8{{end}}

[log]
level = debug

