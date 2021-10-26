package main

import (
	"github.com/jrpalma/pwdhash/config"
	"github.com/jrpalma/pwdhash/logs"
	"github.com/jrpalma/pwdhash/rest"
)

func getLog(conf config.Config) (logs.Logger, error) {
	var err error
	var log logs.Logger

	if conf.LogDestination == logs.FILE {
		log, err = logs.NewFileLogger(conf.LogFile, conf.LogLevel)
		if err != nil {
			return log, err
		}
	}

	log, err = logs.NewStreamLogger(conf.LogDestination, conf.LogLevel)
	return log, err
}

func main() {
	conf := config.Config{}

	err := conf.OpenFile("config.json")
	if err != nil {
		panic(err)
	}

	log, err := getLog(conf)
	if err != nil {
		panic(err)
	}

	server := rest.NewServer(conf, log)
	server.Run()
}
