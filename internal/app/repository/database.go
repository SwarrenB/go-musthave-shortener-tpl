package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/utils"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
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
	dbpool   *pgxpool.Pool
}

// CreateURLRepository implements repository.URLRepository.
func (sqldb *SQLDatabase) CreateURLRepository() (*URLRepositoryState, error) {
	return nil, nil
}

// RestoreURLRepository implements repository.URLRepository.
func (sqldb *SQLDatabase) RestoreURLRepository(m *URLRepositoryState) error {
	panic("unimplemented")
}

func NewDatabase(
	dbOpener SQLDBOpener,
	dataSourceBuilder DataSourceBuilder,
	log zap.Logger,
	driverName string,
	dsn string,
) (*SQLDatabase, error) {
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, errors.New("failed to connect to database")
	}
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
		dbpool:   dbpool,
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

func NewSQLDatabaseConnection(dsn string, log zap.Logger) *SQLDatabase {
	sqldb, err := NewDatabase(sql.Open, PGDataSourceBuilder, log, "pgx", dsn)
	if err != nil {
		log.Fatal("Error creating database")
		return nil
	}
	return sqldb
}

func (sqldb *SQLDatabase) GetURL(shortURL string) (originalURL string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := sqldb.dbpool.QueryRow(ctx, utils.GetURLRegular, shortURL)
	err = row.Scan(&originalURL)
	if err != nil {
		sqldb.log.Error("failed to query url",
			zap.String("short_url", shortURL),
			zap.String("original_url", originalURL),
			zap.Error(err))
	}

	return originalURL, err
}
func (sqldb *SQLDatabase) AddURL(shortURL, originalURL, userID string) (existingURL string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = sqldb.dbpool.QueryRow(ctx, utils.SetURLRegular, shortURL, originalURL, userID).Scan(&existingURL)
	if err != nil || !errors.Is(err, sql.ErrNoRows) {
		sqldb.log.Info("failed to set url",
			zap.String("short_url", shortURL),
			zap.String("original_url", originalURL),
			zap.Error(err))
		return existingURL, err
	} else if errors.Is(err, sql.ErrNoRows) {
		sqldb.log.Info("this url already exists",
			zap.String("short_url", shortURL),
			zap.String("original_url", originalURL),
			zap.Error(err))
		_ = sqldb.database.QueryRow(utils.GetExistingURLRegular, originalURL).Scan(&existingURL)
		return existingURL, err
	}

	sqldb.log.Info("adding to db", zap.String("short_url", shortURL), zap.String("original_url", originalURL))
	return shortURL, nil
}

func (sqldb *SQLDatabase) GetURLByUserID(userID string) ([]Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := sqldb.dbpool.Query(ctx, utils.GetURLsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Record
	var index = 0
	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.ShortURL, &rec.OriginalURL); err != nil {
			return nil, err
		}
		rec.UserID = userID
		rec.ID = index
		results = append(results, rec)
		index++
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
