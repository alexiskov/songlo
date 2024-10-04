package main

import (
	"os"
	"songlib/confreader"
	"songlib/htpsrv"
	"songlib/logger"
	"songlib/psql"
)

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

	if err = htpsrv.New(conf.Server.Port).Start(); err != nil {
		logger.Log.Debug(err.Error())
		return
	}
}
