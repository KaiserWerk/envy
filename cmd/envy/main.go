package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/KaiserWerk/envy/internal/configuration"
	"github.com/KaiserWerk/envy/internal/handler"
	"github.com/KaiserWerk/envy/internal/logging"
	"github.com/KaiserWerk/envy/internal/middleware"
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

	config, err := configuration.FromFile(*configFile)
	if err != nil {
		fmt.Printf("could not read configuration from file '%s': %s\n", *configFile, err.Error())
		fmt.Println("using defaults")
		config = configuration.LoadDefaults()
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

	hd := handler.NewHandler(config, logger)
	_ = hd.LoadVars()

	mwBase := middleware.NewBase(config.AuthKey)

	router := http.NewServeMux()
	router.HandleFunc("/getvar", mwBase.Auth(hd.GetVar))
	router.HandleFunc("/setvar", mwBase.Auth(hd.SetVar))
	router.HandleFunc("/getallvars", mwBase.Auth(hd.GetAllVars))

	srv := http.Server{
		Handler:           router,
		Addr:              config.BindAddress,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    2048,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Errorf("graceful server shutdown failed: %s", err)
		}
		shutdownCancel()
	}()

	logger.Infof("starting up server on %s...", config.BindAddress)
	srv.ListenAndServe()
	logger.Info("goodbye!")
}
