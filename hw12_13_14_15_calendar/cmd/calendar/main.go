package main

import (
	"context"
	"flag"
	sqlstorage "github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/storage/sql"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/app"
	config2 "github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/azicussdu/otus_homework/hw12_13_14_15_calendar/internal/storage/memory"
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

	config, err := config2.NewConfig(configFile)
	if err != nil {
		log.Printf("failed to read config: %v", err)
		os.Exit(1) //nolint:gocritic
	}

	logg := logger.New(config.Logger.Level)

	var storage app.Storage
	switch config.StorageType {
	case "memory":
		storage = memorystorage.New()
	case "database":
		storage, err = sqlstorage.New(getDataSourcePath(config), config.DatabaseConf.MigrationPath)
		if err != nil {
			logg.Error("failed to initialize database storage: " + err.Error())
			os.Exit(1) //nolint:gocritic
		}
	default:
		logg.Error("unknown storage type: " + config.StorageType)
		os.Exit(1) //nolint:gocritic
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, config.ServerConf)

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
		os.Exit(1) //nolint:gocritic
	}
}

func getDataSourcePath(config *config2.Config) string {
	// dsn := "postgres://username:password@localhost:5432/mydatabase"
	return config.DatabaseConf.Type + "://" + config.DatabaseConf.User + ":" +
		config.DatabaseConf.Password + "@" + config.DatabaseConf.Host + ":" +
		strconv.Itoa(config.DatabaseConf.Port) + "/" + config.DatabaseConf.DBName
}
