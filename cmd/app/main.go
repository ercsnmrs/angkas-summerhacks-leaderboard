package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gitlab.angkas.com/avengers/microservice/incentive-service/config"
	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/driver"
	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/kafka"
	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/leaderboard"
	"gitlab.angkas.com/avengers/microservice/incentive-service/logging"
	"gitlab.angkas.com/avengers/microservice/incentive-service/open_loyalty"
	"gitlab.angkas.com/avengers/microservice/incentive-service/server"
	"gitlab.angkas.com/avengers/microservice/incentive-service/storage/redis"
	"gitlab.angkas.com/avengers/microservice/incentive-service/telemetry"
	"gitlab.angkas.com/avengers/microservice/incentive-service/worker"
)

const (
	modeServer = "server"
	modeWorker = "worker"
)

type App struct {
	config   *config.Config
	server   *server.Server
	worker   *worker.Worker
	logger   *slog.Logger
	version  server.Version
	closerFn func() error
}

func (a *App) Setup() error {
	// Init PostgresClient
	// postgresClient, err := postgres.NewClient(a.config.Postgres, a.logger)
	// if err != nil {
	// 	return fmt.Errorf("could not setup postgres: %s", err)
	// }

	auth := &server.JWTAuth{NoVerify: true}
	tsi := telemetry.NewServerInstrumentation(a.config.Telemetry.ServiceName)

	// Init Analytics and Provider
	// Can be Kafka, RabbitMQ, etc.
	ckafkaProducer, err := kafka.NewKafkaProducer(&a.config.KafkaWriter)
	if err != nil {
		return err
	}
	kafkaWriter := kafka.NewClient(a.logger, &a.config.KafkaWriter, ckafkaProducer, time.Now)

	// streamSvc := stream.NewKafkaService(a.logger, kafkaWriter)

	// // Send fake trip data
	// event := trip.GenerateFakeTripEvent()
	// for i := 0; i < 10; i++ {
	// 	streamSvc.SendFakeTripData(context.Background(), event)
	// }

	// Init Provider and Service
	// Can be Redis, Memcached, etc.
	redisClient, err := redis.New(&a.config.Redis, a.logger)
	if err != nil {
		return err
	}
	cacheService := redis.NewCacheService(*redisClient, a.logger)

	// Init Provider and Service
	// Can be any external service OpenLoyalty, TalonOne, etc.
	openLoyaltyClient := open_loyalty.NewOpenLoyaltyClient(
		&a.config.OpenLoyalty,
		cacheService,
		a.logger,
	)
	providerService := open_loyalty.NewProviderSerivce(*openLoyaltyClient, a.logger)

	//svc := foo.NewService(postgresClient, a.logger)
	//service := telemetry.TraceFooService(svc, a.logger)

	driversvc := driver.NewDriverService(cacheService, providerService, a.logger)

	leaderboardsvc := leaderboard.NewLeaderboardService(cacheService, driversvc, a.logger)

	//Generate Fake Drivers
	// driver := faker.GenerateFakeDrivers(15)
	// for _, d := range driver {
	// 	cacheService.RefreshLeaderboard(context.Background(), d)
	// }

	a.server = server.New(a.config.Server, driversvc, leaderboardsvc, auth, tsi, a.version, a.logger)

	//tiersvc.RefreshTier(context.Background())

	// Check if we need to use same implementation.
	// fj := fakejob.New(time.Second)

	a.worker = worker.New(kafkaWriter, 1, a.logger)
	a.worker.Use(worker.LoggingMiddleware(a.logger), telemetry.TraceWorker)
	a.worker.HandleFunc("trips", worker.ConsumeTripCompleted(leaderboardsvc))

	// refreshTierWeekly, err := worker.NewSchedule(
	// 	// a.conf.ConfigResetClock.Format(time.Kitchen),
	// 	"",
	// 	time.Hour*24,
	// 	func(ctx context.Context) error { return tiersvc.RefreshTier(ctx) },
	// )
	// if err != nil {
	// 	return fmt.Errorf("could not setup refreshing tier weekly: %s", err)
	// }
	// a.worker.SetSchedule(refreshTierWeekly)

	a.closerFn = func() error {
		// if err = postgresClient.Close(); err != nil {
		// 	return fmt.Errorf("could not close postgres: %s", err)
		// }
		if err = cacheService.Client.Close(); err != nil {
			return fmt.Errorf("could not close postgres: %s", err)
		}
		if err = a.server.Close(); err != nil {
			return fmt.Errorf("could not close server: %s", err)
		}
		return nil
	}
	return nil
}

func (a *App) Run(mode string) error {
	switch mode {
	case modeServer:
		return appRunner(a.server)
	case modeWorker:
		return appRunner(a.worker)
	default:
		return fmt.Errorf("app mode not supported: %s", mode)
	}
}

func main() {
	log := logging.Default()

	mode, err := appMode()
	if err != nil {
		log.Error("app mode", "err", err)
		return
	}

	log.Info("app: loading config...")
	conf, err := config.LoadDefault()
	if err != nil {
		log.Error("could not load config", "err", err)
		return
	}

	log.Info(fmt.Sprintf("log level set to %s", strings.ToUpper(conf.Logging.Level)))
	log, err = logging.New(conf.Logging)
	if err != nil {
		log.Error("could not create logging", "err", err)
		return
	}

	version := buildVer()
	log.Info(fmt.Sprintf("telemetry enabled:%v url:%s", conf.Telemetry.Enabled, conf.Telemetry.CollectorURL))
	shutdown, err := telemetry.InitProvider(conf.Telemetry, mode, version.Tag)
	if err != nil {
		log.Error("could not init telemetry provider", "err", err)
		return
	}
	defer func() {
		if err = shutdown(context.Background()); err != nil {
			log.Error("failed to shutdown TracerProvider: %w", err)
		}
	}()
	log = telemetry.TraceLogger(log)

	app := &App{config: conf, logger: log, version: version}
	if err = app.Setup(); err != nil {
		log.Error("could not setup app", "err", err)
		return
	}
	log.Info("running", "mode", mode, "version", version)
	if err = app.Run(mode); err != nil {
		log.Error("could not run app", "err", err)
		return
	}
	if err = app.closerFn(); err != nil {
		log.Error("error occurred when closing app", "err", err)
	}
	log.Info("app: exited!")
}

// appMode returns application mode base on arguments.
func appMode() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("app mode required")
	}
	return strings.ToLower(os.Args[1]), nil
}

func appRunner(app runner) error {
	done := make(chan error, 1)
	// Waits for CTRL-C or os SIGINT for server shutdown.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		done <- app.Stop()
	}()

	if err := app.Run(); err != nil {
		return err
	}
	return <-done
}

type runner interface {
	Run() error
	Stop() error
}

// version data will set during build time using go build -ldflags.
var vTag, vCommit, vBuilt string

func buildVer() server.Version {
	ts, _ := strconv.Atoi(vBuilt)
	bt := time.Unix(int64(ts), 0)
	return server.Version{
		Tag:    vTag,
		Commit: vCommit,
		Built:  bt,
	}
}
