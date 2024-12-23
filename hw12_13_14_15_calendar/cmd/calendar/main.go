package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/app"                          //nolint:depguard
	"github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/config"                       //nolint:depguard
	"github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/logger"                       //nolint:depguard
	internalhttp "github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/server/http"     //nolint:depguard
	memorystorage "github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/storage/memory" //nolint:depguard
	sqlstorage "github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"       //nolint:depguard
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		log.Printf("failed to read config: %v", err)
		os.Exit(1)
	}

	logg := logger.New(cfg.Logger.Level)

	var storage app.Storage
	switch cfg.StorageType {
	case "memory":
		storage = memorystorage.New()
	case "database":
		storage, err = sqlstorage.New(cfg.DatabaseConf)
		if err != nil {
			logg.Error("failed to initialize database storage: " + err.Error())
			os.Exit(1)
		}
	default:
		logg.Error("unknown storage type: " + cfg.StorageType)
		os.Exit(1)
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, cfg.ServerConf)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*3)
		defer shutdownCancel()

		if err = server.Stop(shutdownCtx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err = server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
	}
}
