package stores

import (
	"context"
	"fmt"

	"github.com/tuantran1810/go-di-template/internal/models"
)

const DefaultLimit = 100

type GenericStore[T any] struct {
	*Repository
}

func NewGenericStore[T any](repository *Repository) *GenericStore[T] {
	return &GenericStore[T]{Repository: repository}
}

func (s *GenericStore[T]) Ping(ctx context.Context) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var entity T
	dbtx := s.GetTransaction(nil).WithContext(timeoutCtx)
	if err := dbtx.Limit(1).Select("id").Find(&entity).Error; err != nil {
		return GenerateError("failed to ping database", err)
	}

	return nil
}

func (s *GenericStore[T]) Create(
	ctx context.Context,
	tx models.Transaction,
	entity *T,
) (*T, error) {
	if entity == nil {
		return nil, fmt.Errorf("%w - input entity is nil", models.ErrInvalid)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	dbtx := s.GetTransaction(tx).WithContext(timeoutCtx)
	if err := dbtx.Create(entity).Error; err != nil {
		return nil, GenerateError("failed to create entity", err)
	}

	return entity, nil
}

func (s *GenericStore[T]) CreateMany(
	ctx context.Context,
	tx models.Transaction,
	entities []T,
) ([]T, error) {
	if len(entities) == 0 {
		return nil, fmt.Errorf("%w - input entities is empty", models.ErrInvalid)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	dbtx := s.GetTransaction(tx).WithContext(timeoutCtx)
	if err := dbtx.Create(entities).Error; err != nil {
		return nil, GenerateError("failed to create entities", err)
	}

	return entities, nil
}

func (s *GenericStore[T]) Get(
	ctx context.Context,
	tx models.Transaction,
	id uint,
) (*T, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	dbtx := s.GetTransaction(tx).WithContext(timeoutCtx)
	var entity T
	if err := dbtx.First(&entity, id).Error; err != nil {
		return nil, GenerateError("failed to get entity", err)
	}

	return &entity, nil
}

func (s *GenericStore[T]) GetMany(
	ctx context.Context,
	tx models.Transaction,
	ids []uint,
) ([]*T, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("%w - input ids is empty", models.ErrInvalid)
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	dbtx := s.GetTransaction(tx).WithContext(timeoutCtx)
	var entities []*T
	if err := dbtx.Find(&entities, ids).Error; err != nil {
		return nil, GenerateError("failed to get entities", err)
	}

	return entities, nil
}

func (s *GenericStore[T]) GetByCriterias(
	ctx context.Context,
	tx models.Transaction,
	fields []string,
	criterias map[string]any,
	orderBys []string,
) (*T, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var entity T
	dbtx := s.GetTransaction(tx).WithContext(timeoutCtx)
	if len(fields) > 0 {
		dbtx = dbtx.Select(fields)
	}
	for k, v := range criterias {
		dbtx = dbtx.Where(fmt.Sprintf("`%s` = ?", k), v)
	}
	for _, order := range orderBys {
		dbtx = dbtx.Order(order)
	}

	if err := dbtx.First(&entity).Error; err != nil {
		return nil, GenerateError("failed to get entity", err)
	}

	return &entity, nil
}

func (s *GenericStore[T]) GetManyByCriterias(
	ctx context.Context,
	tx models.Transaction,
	fields []string,
	criterias map[string]any,
	orderBys []string,
	offset int,
	limit int,
) ([]*T, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	dbtx := s.GetTransaction(tx).WithContext(timeoutCtx)

	for k, v := range criterias {
		dbtx = dbtx.Where(fmt.Sprintf("`%s` = ?", k), v)
	}
	for _, order := range orderBys {
		dbtx = dbtx.Order(order)
	}

	if limit <= 0 {
		limit = DefaultLimit
	}
	var entities []*T
	if err := dbtx.
		Offset(offset).
		Limit(limit).
		Select(fields).
		Find(&entities).
		Error; err != nil {
		return nil, GenerateError("failed to get entity", err)
	}

	return entities, nil
}

func (s *GenericStore[T]) Count(
	ctx context.Context,
	tx models.Transaction,
	criterias map[string]any,
) (int64, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	dbtx := s.GetTransaction(tx).WithContext(timeoutCtx)

	for k, v := range criterias {
		dbtx = dbtx.Where(fmt.Sprintf("`%s` = ?", k), v)
	}

	var entity T
	var cnt int64
	if err := dbtx.Model(&entity).Count(&cnt).Error; err != nil {
		return 0, GenerateError("failed to count", err)
	}

	return cnt, nil
}

func (s *GenericStore[T]) Update(
	ctx context.Context,
	tx models.Transaction,
	entity *T,
) error {
	if entity == nil {
		return fmt.Errorf("%w - input entity is nil", models.ErrInvalid)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	dbtx := s.GetTransaction(tx).WithContext(timeoutCtx)
	dbtx = dbtx.Updates(entity)
	if err := dbtx.Error; err != nil {
		return GenerateError("failed to update entity", err)
	}
	if dbtx.RowsAffected == 0 {
		return fmt.Errorf("%w - no rows affected", models.ErrNotFound)
	}

	return nil
}

func (s *GenericStore[T]) Delete(
	ctx context.Context,
	tx models.Transaction,
	permanent bool,
	id uint,
) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	dbtx := s.GetTransaction(tx).WithContext(timeoutCtx)
	if permanent {
		dbtx = dbtx.Unscoped()
	}

	var entity T
	if err := dbtx.
		Model(&entity).
		Delete("id = ?", id).Error; err != nil {
		return GenerateError("failed to delete entity", err)
	}

	return nil
}

func (s *GenericStore[T]) DeleteMany(
	ctx context.Context,
	tx models.Transaction,
	permanent bool,
	ids []uint,
) (int64, error) {
	if len(ids) == 0 {
		return 0, fmt.Errorf("%w - input ids is empty", models.ErrInvalid)
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	dbtx := s.GetTransaction(tx).WithContext(timeoutCtx)
	if permanent {
		dbtx = dbtx.Unscoped()
	}

	var entity T
	dbtx = dbtx.Model(&entity).Delete("id in (?)", ids)
	if err := dbtx.Error; err != nil {
		return 0, GenerateError("failed to delete entities", err)
	}

	return dbtx.RowsAffected, nil
}
