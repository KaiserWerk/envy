package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/KaiserWerk/envy/internal/configuration"
	"github.com/KaiserWerk/envy/internal/logging"
)

// version info
var (
	Version = "DEV"
	Date    = "2022-01-01 00:00:00"
	Commit  = "00000000"
)

// flags
var (
	configFile = flag.String("cfg", "app.yaml", "The configuration file to use")
	logFile    = flag.String("logfile", "envy.log", "The log file to write")
)

func main() {
	flag.Parse()

	_, err := configuration.FromFile(*configFile)
	if err != nil {
		fmt.Printf("could not read configuration from file '%s': %s\n", *configFile, err.Error())
		os.Exit(-1)
	}

	fh, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("could not open log file for writing: " + err.Error())
		os.Exit(-1)
	}
	defer fh.Close()

	logger := logging.New(fh)
	logger.Infof("Envy server; Version: %s; Version date: %s; Git Commit: %s", Version, Date, Commit)
	logger.Infof("Used configuration file: '%s'", *configFile)

}
