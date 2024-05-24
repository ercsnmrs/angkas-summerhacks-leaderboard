package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
	"gitlab.angkas.com/avengers/microservice/incentive-service/internal/kafka"
	"gitlab.angkas.com/avengers/microservice/incentive-service/logging"
	"gitlab.angkas.com/avengers/microservice/incentive-service/open_loyalty"
	"gitlab.angkas.com/avengers/microservice/incentive-service/server"
	"gitlab.angkas.com/avengers/microservice/incentive-service/storage/postgres"
	"gitlab.angkas.com/avengers/microservice/incentive-service/storage/redis"
	"gitlab.angkas.com/avengers/microservice/incentive-service/telemetry"
)

const DefaultFile = ".env"

// Config represents application configuration.
type Config struct {
	Server                       server.Config
	WorkerQueueSize              int
	Logging                      logging.Config
	Telemetry                    telemetry.Config
	GoogleApplicationCredentials string
	Postgres                     postgres.Config
	Redis                        redis.Config
	OpenLoyalty                  open_loyalty.Config
	KafkaWriter                  kafka.WriterConfig
}

// Load loads config from environment variables and file.
func Load(file string) (*Config, error) {
	viper.SetConfigFile(file)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Set Default Values for Kafka
	viper.SetDefault("BUDGET_MONITORING_INTERVAL", 10)
	viper.SetDefault("KAFKA_PRODUCER_BATCH_SIZE", 100000)
	viper.SetDefault("KAFKA_PRODUCER_LINGER_MS", 10)
	viper.SetDefault("KAFKA_PRODUCER_COMPRESSION_TYPE", "lz4")
	viper.SetDefault("KAFKA_PRODUCER_ACKS", "all")
	viper.SetDefault("KAFKA_PRODUCER_TIMEOUT", 30000)

	if err := viper.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	c := &Config{
		Server: server.Config{
			Addr:         viper.GetString("SERVER_ADDR"),
			ReadTimeout:  viper.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout: viper.GetDuration("SERVER_WRITE_TIMEOUT"),
		},
		WorkerQueueSize: viper.GetInt("WORKER_QUEUE_SIZE"),
		Logging: logging.Config{
			Level: viper.GetString("LOGGING_LEVEL"),
		},
		Telemetry: telemetry.Config{
			Enabled:      viper.GetBool("TELEMETRY_ENABLED"),
			CollectorURL: viper.GetString("TELEMETRY_COLLECTOR_URL"),
			ServiceName:  viper.GetString("TELEMETRY_SERVICE_NAME"),
			Env:          viper.GetString("TELEMETRY_ENV"),
		},
		Postgres: postgres.Config{
			URL:             viper.GetString("POSTGRES_URL"),
			MaxConns:        viper.GetInt("POSTGRES_MAX_CONNS"),
			MaxConnIdleTime: viper.GetDuration("POSTGRES_MAX_IDLE_TIME"),
			MaxConnLifetime: viper.GetDuration("POSTGRES_MAX_LIFE_TIME"),
		},
		Redis: redis.Config{
			Host:     viper.GetString("REDIS_URL"),
			Password: viper.GetString("REDIS_PASSWORD"),
			Port:     viper.GetString("REDIS_PORT"),
		},
		OpenLoyalty: open_loyalty.Config{
			URL:      viper.GetString("OPENLOYALTY_URL"),
			Username: viper.GetString("OPENLOYALTY_USERNAME"),
			Password: viper.GetString("OPENLOYALTY_PASSWORD"),
			StoreID:  viper.GetString("OPENLOYALTY_STORE_ID"),
		},
		KafkaWriter: kafka.WriterConfig{
			Servers:          viper.GetString("KAFKA_SERVERS"),
			SecurityProtocol: viper.GetString("KAFKA_SECURITY_PROTOCOL"),
			SSLMechanism:     viper.GetString("KAFKA_SSL_MECHANISM"),
			Username:         viper.GetString("KAFKA_USERNAME"),
			Password:         viper.GetString("KAFKA_PASSWORD"),
			Topic:            viper.GetString("KAFKA_WRITER_EVENTS_TOPIC"),
			FlushTimeout:     int(viper.GetDuration("KAFKA_WRITER_EVENTS_FLUSH_TIMEOUT")),
			BatchSize:        viper.GetInt("KAFKA_PRODUCER_BATCH_SIZE"),
			LingerMs:         viper.GetInt("KAFKA_PRODUCER_LINGER_MS"),
			CompressionType:  viper.GetString("KAFKA_PRODUCER_COMPRESSION_TYPE"),
			Acks:             viper.GetString("KAFKA_PRODUCER_ACKS"),
			ProducerTimeout:  viper.GetInt("KAFKA_PRODUCER_TIMEOUT"),
		},
		GoogleApplicationCredentials: viper.GetString("GOOGLE_APPLICATION_CREDENTIALS"),
	}
	return c, nil
}

// LoadDefault loads config from environment variables and .env file.
func LoadDefault() (*Config, error) {
	return Load(DefaultFile)
}
