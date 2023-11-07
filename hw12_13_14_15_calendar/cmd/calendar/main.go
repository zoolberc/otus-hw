package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/internal/server/http"
	st "github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/zoolberc/otus-hw/hw12_13_14_15_calendar/internal/storage/sql"
)

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	logfile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.Println("error when initializing log file: ", err)
		os.Exit(1)
	}
	defer logfile.Close()

	log := logger.Setup(config.LogLevel, logfile)

	storage := setStorage(config)

	calendar := app.New(log, storage)

	server := internalhttp.NewServer(log, config.Host, config.Port, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Error("failed to stop http server: " + err.Error())
		}
		log.Info("graceful shutdown complete")
	}()

	log.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		log.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func setStorage(config Config) st.Storage {
	if config.StorageType == "sql" {
		return sqlstorage.New(config.DataBaseConf)
	}
	return memorystorage.New()
}
