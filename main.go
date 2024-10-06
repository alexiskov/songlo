package main

import (
	"os"
	"songlib/confreader"
	"songlib/htpsrv"
	"songlib/logger"
	"songlib/psql"
	"strconv"

	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

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
