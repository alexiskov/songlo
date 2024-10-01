package main

import (
	"log"
	"songlib/confreader"
	"songlib/htpsrv"
	"songlib/psql"
)

func main() {
	conf, err := confreader.LoadConfig()
	if err != nil {
		log.Println(err)
		return
	}
	err = psql.Init(conf.DMS.Host, conf.DMS.Username, conf.DMS.Password, conf.DMS.DBname, conf.DMS.Port)
	if err != nil {
		log.Println(err)
		return
	}

	if err = htpsrv.New(conf.Server.Port).Start(); err != nil {
		log.Println(err)
		return
	}
}
