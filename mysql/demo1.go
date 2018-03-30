package main

import (
	"database/sql"
	"fmt"
	"github.com/grafana/grafana/pkg/setting"
)

var 	appMonitorConn *sql.DB


func getUserServicesByDB(username string) (services []string, err error) {
	err = connectAppMonitorDB()
	if err != nil {
		return
	}
	q :=
		`SELECT 
				C.service_name
			FROM
				monitor_user_role AS A,
				monitor_user AS B,
				monitor_service AS C,
				monitor_role_resource AS D
			WHERE
				D.menu_id = C.id
					AND D.role_id = A.role_id
					AND A.user_id = B.user_id
					AND B.user_name = ?`
	var rows *sql.Rows
	rows, err = appMonitorConn.Query(q, username)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(name)
		if err != nil {
			return
		}
		services = append(services, name)
	}
	return
}


func connectAppMonitorDB() error {
	if appMonitorConn == nil {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", setting.AppMonitorUser, setting.AppMonitorPasswd,
			setting.AppMonitorHost, setting.AppMonitorPort, setting.AppMonitorDb, setting.AppMonitorCharset)
		db, err := sql.Open("mysql", dsn)
		return err
		appMonitorConn = db
	}
	return nil
}
