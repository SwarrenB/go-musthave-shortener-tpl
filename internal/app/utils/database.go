package utils

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type DBConfig struct {
	Host string
	User string
	DB   string
}

type (
	SQLDBOpener       func(driverName, dataSourceName string) (*sql.DB, error)
	DataSourceBuilder func(dsn string) (string, *DBConfig, error)
)

type SQLDatabase struct {
	database *sql.DB
	dbConfig *DBConfig
	log      zap.Logger
}

func NewDatabase(
	dbOpener SQLDBOpener,
	dataSourceBuilder DataSourceBuilder,
	log zap.Logger,
	driverName string,
	dsn string,
) (*SQLDatabase, error) {
	dataSourceName, dbConfig, err := dataSourceBuilder(dsn)
	if err != nil {
		return nil, err
	}

	sqldb, err := dbOpener(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	return &SQLDatabase{
		database: sqldb,
		dbConfig: dbConfig,
		log:      log,
	}, nil
}

func (sqldb *SQLDatabase) Ping() error {
	if err := sqldb.database.Ping(); err != nil {
		return err
	}

	sqldb.log.Info(
		"Database has been started",
		zap.String("host", sqldb.dbConfig.Host),
		zap.String("database", sqldb.dbConfig.DB),
		zap.String("user", sqldb.dbConfig.User),
	)

	return nil
}

func (sqldb *SQLDatabase) Close() error {
	if err := sqldb.database.Close(); err != nil {
		return err
	}

	return nil
}

func PGDataSourceBuilder(dsn string) (string, *DBConfig, error) {
	config, err := pgconn.ParseConfig(dsn)
	if err != nil {
		return ``, nil, err
	}

	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.User, config.Password, config.Database)

	return dataSourceName, &DBConfig{
		Host: config.Host,
		User: config.User,
		DB:   config.Database,
	}, nil
}

func NewPG(dsn string, log zap.Logger) (*SQLDatabase, error) {
	sqldb, err := NewDatabase(sql.Open, PGDataSourceBuilder, log, "pgx", dsn)
	if err != nil {
		return nil, err
	}
	return sqldb, nil
}
