package stores

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/tuantran1810/go-di-template/internal/entities"
	"gorm.io/driver/sqlite"
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
	DatabasePath string
}

type Repository struct {
	mutex sync.RWMutex
	db    *gorm.DB
}

func NewRepository(cfg RepositoryConfig) (*Repository, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DatabasePath), &gorm.Config{}) //nolint: varnamelen
	if err != nil {
		return nil, fmt.Errorf("%w - failed to open database: %w", entities.ErrDatabase, err)
	}

	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, fmt.Errorf("%w - failed to enable foreign keys: %w", entities.ErrDatabase, err)
	}

	if err := db.Exec("PRAGMA journal_mode = WAL").Error; err != nil {
		return nil, fmt.Errorf("%w - failed to enable WAL mode: %w", entities.ErrDatabase, err)
	}

	return &Repository{db: db}, nil
}

func MustNewRepository(cfg RepositoryConfig) *Repository {
	repo, err := NewRepository(cfg)
	if err != nil {
		panic(err)
	}

	return repo
}

func (r *Repository) Start(_ context.Context) error {
	return nil
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

	return txImpl.(*gorm.DB) //nolint: forcetypeassert
}

func (r *Repository) RunTx(ctx context.Context, data any, funcs ...entities.DBTxHandleFunc) (any, error) {
	if len(funcs) == 0 {
		return data, fmt.Errorf("%w - input no handler function", entities.ErrInternal)
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txKeeper := entities.NewGormTransaction(tx)

		for _, f := range funcs {
			if f != nil {
				outData, cont, ferr := f(ctx, txKeeper, data)
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
