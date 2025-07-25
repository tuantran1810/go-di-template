package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tuantran1810/go-di-template/uberfx/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultTimeout = 20 * time.Second

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
		return models.ErrCanceled
	}

	if isInvalidInputError(err) {
		return models.ErrInvalid
	}

	if isNotFoundError(err) {
		return models.ErrNotFound
	}

	return models.ErrDatabase
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

	eligibleErr := errors.Is(err, models.ErrCanceled) ||
		errors.Is(err, models.ErrInvalid) ||
		errors.Is(err, models.ErrNotFound) ||
		errors.Is(err, models.ErrDatabase)
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
		return fmt.Errorf("%w - failed to open database: %w", models.ErrDatabase, err)
	}

	dbInstance, err := db.DB()
	if err != nil {
		return fmt.Errorf("%w - failed to get database instance: %w", models.ErrDatabase, err)
	}

	dbInstance.SetMaxOpenConns(int(r.RepositoryConfig.MaxOpenConns))
	dbInstance.SetMaxIdleConns(int(r.RepositoryConfig.MaxIdleConns))
	dbInstance.SetConnMaxLifetime(time.Duration(r.RepositoryConfig.ConnMaxLifeTimeSeconds) * time.Second)

	if err := dbInstance.Ping(); err != nil {
		return fmt.Errorf("%w - failed to ping database: %w", models.ErrDatabase, err)
	}

	r.db = db
	return r.Check(ctx)
}

func (r *Repository) Stop(_ context.Context) error {
	db, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("%w - failed to get database connection: %w", models.ErrDatabase, err)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("%w - failed to close database connection: %w", models.ErrDatabase, err)
	}

	return nil
}

func (r *Repository) Check(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	if err := r.db.WithContext(timeoutCtx).Exec("SELECT 1").Error; err != nil {
		return GenerateError("failed to ping database", err)
	}

	return nil
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

func (r *Repository) GetTransaction(tx models.Transaction) *gorm.DB {
	if tx == nil {
		return r.db
	}

	txImpl := tx.GetTransaction()
	if txImpl == nil {
		return r.db
	}

	return txImpl.(*gorm.DB) //nolint: forcetypeassert
}

func (r *Repository) RunTx(ctx context.Context, data any, funcs ...models.DBTxHandleFunc) (any, error) {
	if len(funcs) == 0 {
		return data, fmt.Errorf("%w - input no handler function", models.ErrInternal)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	err := r.db.WithContext(timeoutCtx).Transaction(func(tx *gorm.DB) error {
		txKeeper := models.NewGormTransaction(tx)

		for _, f := range funcs {
			if f != nil {
				outData, cont, ferr := f(timeoutCtx, txKeeper, data)
				data = outData
				if ferr != nil {
					return ferr
				}

				if !cont {
					break
				}
			}
		}

		return nil
	})

	return data, handleTransactionError(err)
}
