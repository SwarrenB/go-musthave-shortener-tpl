package utils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
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
		return nil, errors.New("error parsing database DSN")
	}

	sqldb, err := dbOpener(driverName, dataSourceName)
	if err != nil {
		return nil, errors.New("error opening database")
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

func Init(dsn string, log zap.Logger) *SQLDatabase {
	sqldb, err := NewDatabase(sql.Open, PGDataSourceBuilder, log, "pgx", dsn)
	if err != nil {
		log.Fatal("Error creating database")
		return nil
	} else {
		sqldb.CreateTables(log)
	}
	return sqldb
}

func (sqldb *SQLDatabase) CreateTables(logger zap.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := sqldb.database.ExecContext(ctx, createTableQuery); err != nil {
		logger.Fatal("Failed to create tables",
			zap.Error(err),
			zap.String("Query", createTableQuery),
		)
	}
}

func (sqldb *SQLDatabase) GetFromDB(shortURL string) (originalURL string, ok bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := sqldb.database.QueryRowContext(ctx, getURLRegular, shortURL)
	err := row.Scan(&originalURL)
	if err != nil {
		sqldb.log.Error("failed to query url",
			zap.String("short_url", shortURL),
			zap.String("original_url", originalURL),
			zap.Error(err))
	}

	return originalURL, row != nil
}
func (sqldb *SQLDatabase) AddToDB(shortURL, originalURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := sqldb.database.ExecContext(ctx, setURLRegular, shortURL, originalURL)
	if err != nil {
		sqldb.log.Error("failed to set url",
			zap.String("short_url", shortURL),
			zap.String("original_url", originalURL),
			zap.Error(err))
	}
}
