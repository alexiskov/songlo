package main

import (
	"os"
	"songlib/confreader"
	"songlib/htpsrv"
	"songlib/logger"
	"songlib/psql"
	"strconv"

	_ "songlib/docs"
)

//	@title			songlibs
//	@version		0.0.2
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	personal
//	@license.url	http://www.youtube.com

//	@host		localhost:8080
//	@BasePath	/

func main() {

	logger.New(os.Stdout)
	logger.Log.Info("logger is ready!")

	conf, err := confreader.LoadConfig()
	if err != nil {
		logger.Log.Debug(err.Error())
		return
	}
	logger.Log.Info("configs loaded...")

	err = psql.Init(conf.DMS.Host, conf.DMS.Username, conf.DMS.Password, conf.DMS.DBname, conf.DMS.Port)
	if err != nil {
		logger.Log.Debug(err.Error())
		return
	}
	logger.Log.Info("dataBase initialized")

	if err = htpsrv.Start(strconv.Itoa(conf.Server.Port)); err != nil {
		logger.Log.Debug(err.Error())
		return
	}
}
