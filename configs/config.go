package configs

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

type Conf struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
	RabbitMQURL       string `mapstructure:"RABBITMQ_URL"`
}

func LoadConfig(path string) (*Conf, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")
	viper.AddConfigPath("cmd/ordersystem")
	_ = viper.ReadInConfig() // ignora erro; env vars servem como fallback
	var cfg Conf
	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return &cfg, err
}

func EnsureDatabaseAndTable(cfg *Conf) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort)
	db, err := sql.Open(cfg.DBDriver, dsn)
	if err != nil {
		return fmt.Errorf("abrir conexão mysql: %w", err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + cfg.DBName); err != nil {
		return fmt.Errorf("criar banco %q: %w", cfg.DBName, err)
	}

	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db2, err := sql.Open(cfg.DBDriver, dsnWithDB)
	if err != nil {
		return fmt.Errorf("abrir conexão com banco: %w", err)
	}
	defer db2.Close()

	_, err = db2.Exec(`CREATE TABLE IF NOT EXISTS orders (
		id VARCHAR(255) PRIMARY KEY,
		price FLOAT NOT NULL,
		tax FLOAT NOT NULL,
		final_price FLOAT NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("criar tabela orders: %w", err)
	}

	return nil
}
