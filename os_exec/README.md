# 从Mysql查询需要监控的主机和端口， 分别调用系统执行命令，获取agent的信息，再把监控信息推送到influxdb

go run agent_monitor.go -c config.json
