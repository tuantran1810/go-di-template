package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tuantran1810/go-di-template/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func isInvalidInputError(err error) bool {
	if err == nil {
		return false
	}

	if strings.Contains(err.Error(), "constraint") {
		return true
	}

	return errors.Is(err, gorm.ErrInvalidData) ||
		errors.Is(err, gorm.ErrInvalidField) ||
		errors.Is(err, gorm.ErrInvalidValue) ||
		errors.Is(err, gorm.ErrInvalidValueOfLength)
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, gorm.ErrRecordNotFound)
}

func isCanceledError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.Canceled) {
		return true
	}

	errString := err.Error()

	return strings.Contains(errString, "operation was canceled")
}

func getEntityError(err error) error {
	if err == nil {
		return nil
	}

	if isCanceledError(err) {
		return entities.ErrCanceled
	}

	if isInvalidInputError(err) {
		return entities.ErrInvalid
	}

	if isNotFoundError(err) {
		return entities.ErrNotFound
	}

	return entities.ErrDatabase
}

func GenerateError(errStr string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w - %s, err: %w", getEntityError(err), errStr, err)
}

func handleTransactionError(err error) error {
	if err == nil {
		return nil
	}

	eligibleErr := errors.Is(err, entities.ErrCanceled) ||
		errors.Is(err, entities.ErrInvalid) ||
		errors.Is(err, entities.ErrNotFound) ||
		errors.Is(err, entities.ErrDatabase)
	if eligibleErr {
		return err
	}

	return GenerateError("transaction error", err)
}

type RepositoryConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	SSLMode  *string
	Timezone *string
	Params   map[string]string

	Timeout                time.Duration
	MaxOpenConns           uint32
	MaxIdleConns           uint32
	ConnMaxLifeTimeSeconds uint32
}

func (cfg RepositoryConfig) DSN() string {
	basic := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database,
	)

	parts := make([]string, 0)
	parts = append(parts, basic)
	if cfg.SSLMode != nil && *cfg.SSLMode != "" {
		parts = append(parts, fmt.Sprintf("sslmode=%s", *cfg.SSLMode))
	}
	if cfg.Timezone != nil && *cfg.Timezone != "" {
		parts = append(parts, fmt.Sprintf("TimeZone=%s", *cfg.Timezone))
	}
	for k, v := range cfg.Params {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(parts, " ")
}

type Repository struct {
	RepositoryConfig
	db *gorm.DB
}

func MustNewRepository(cfg RepositoryConfig) *Repository {
	return &Repository{
		RepositoryConfig: cfg,
	}
}

func (r *Repository) Start(ctx context.Context) error {
	dsn := r.RepositoryConfig.DSN()
	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		return fmt.Errorf("%w - failed to open database: %w", entities.ErrDatabase, err)
	}

	dbInstance, err := db.DB()
	if err != nil {
		return fmt.Errorf("%w - failed to get database instance: %w", entities.ErrDatabase, err)
	}

	dbInstance.SetMaxOpenConns(int(r.RepositoryConfig.MaxOpenConns))
	dbInstance.SetMaxIdleConns(int(r.RepositoryConfig.MaxIdleConns))
	dbInstance.SetConnMaxLifetime(time.Duration(r.RepositoryConfig.ConnMaxLifeTimeSeconds) * time.Second)

	if err := dbInstance.Ping(); err != nil {
		return fmt.Errorf("%w - failed to ping database: %w", entities.ErrDatabase, err)
	}

	r.db = db
	return r.Check(ctx)
}

func (r *Repository) Stop(_ context.Context) error {
	db, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("%w - failed to get database connection: %w", entities.ErrDatabase, err)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("%w - failed to close database connection: %w", entities.ErrDatabase, err)
	}

	return nil
}

func (r *Repository) Check(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
		return GenerateError("failed to ping database", err)
	}

	return nil
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

func (r *Repository) GetTransaction(tx entities.Transaction) *gorm.DB {
	if tx == nil {
		return r.db
	}

	txImpl := tx.GetTransaction()
	if txImpl == nil {
		return r.db
	}

	return txImpl.(*gorm.DB)
}

func (r *Repository) RunTx(ctx context.Context, funcs ...entities.DBTxHandleFunc) error {
	if len(funcs) == 0 {
		return fmt.Errorf("%w - input no handler function", entities.ErrInternal)
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txKeeper := entities.NewGormTransaction(tx)

		for _, f := range funcs {
			if f != nil {
				ferr := f(ctx, txKeeper)
				if ferr != nil {
					return ferr
				}
			}
		}

		return nil
	})

	return handleTransactionError(err)
}
