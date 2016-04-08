// main
package main

import (
	"apiserver/liveapi"
	"apiserver/service"
	"os"

	"github.com/futurez/litego/config"
	"github.com/futurez/litego/logger"
)

func main() {
	logger.Info("start apiserver!!!")

	iniconf, err := config.NewConfig(config.IniProtocol, "config.ini")
	if err != nil {
		logger.Error("main : ", err)
		os.Exit(1)
	}

	//init mysql
	mysql_addr := iniconf.GetString("mysql.address")
	mysql_dbname := iniconf.GetString("mysql.dbname")
	mysql_account := iniconf.GetString("mysql.account")
	mysql_passwd := iniconf.GetString("mysql.password")
	service.InitMysqlDbPool(mysql_addr, mysql_dbname, mysql_account, mysql_passwd)

	//start http server.
	//server_ip := iniconf.GetString("apiserver.server_ip", "*")
	http_port, _ := iniconf.GetInt("apiserver.http_port", 80)
	liveapi.StartHttpServer("", http_port)

	logger.StartAsyncSave()

	waiting := make(chan bool)
	logger.Info("exit ", <-waiting)
}
