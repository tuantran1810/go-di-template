package postgres

import (
	"context"
	"fmt"

	"github.com/tuantran1810/go-di-template/internal/entities"
)

const DefaultLimit = 100

type GenericRepository[T, E any] struct {
	*Repository
	transformer *entities.ExtendedDataTransformer[T, E]
}

func NewGenericRepository[T, E any](
	repository *Repository,
	transformer *entities.ExtendedDataTransformer[T, E],
) *GenericRepository[T, E] {
	return &GenericRepository[T, E]{
		Repository:  repository,
		transformer: transformer,
	}
}

func (s *GenericRepository[T, E]) Ping(ctx context.Context) error {
	var entity T
	dbtx := s.GetTransaction(nil).WithContext(ctx)
	if err := dbtx.Limit(1).Select("id").Find(&entity).Error; err != nil {
		return GenerateError("failed to ping database", err)
	}

	return nil
}

func (s *GenericRepository[T, E]) AutoMigrate(ctx context.Context) error {
	var entity T
	return s.db.WithContext(ctx).AutoMigrate(&entity)
}

func (s *GenericRepository[T, E]) Create(
	ctx context.Context,
	tx entities.Transaction,
	entity *E,
) (*E, error) {
	if entity == nil {
		return nil, fmt.Errorf("%w - input entity is nil", entities.ErrInvalid)
	}

	data, err := s.transformer.FromEntity(entity)
	if err != nil {
		return nil, err
	}
	dbtx := s.GetTransaction(tx).WithContext(ctx)
	if err := dbtx.Create(data).Error; err != nil {
		return nil, GenerateError("failed to create data", err)
	}

	return s.transformer.ToEntity(data)
}

func (s *GenericRepository[T, E]) CreateMany(
	ctx context.Context,
	tx entities.Transaction,
	entityArray []E,
) ([]E, error) {
	if len(entityArray) == 0 {
		return nil, fmt.Errorf("%w - input entities is empty", entities.ErrInvalid)
	}

	dataArray, err := s.transformer.FromEntityArray_I2I(entityArray)
	if err != nil {
		return nil, err
	}

	dbtx := s.GetTransaction(tx).WithContext(ctx)
	if err := dbtx.Create(dataArray).Error; err != nil {
		return nil, GenerateError("failed to create data records", err)
	}

	return s.transformer.ToEntityArray_I2I(dataArray)
}

func (s *GenericRepository[T, E]) Get(
	ctx context.Context,
	tx entities.Transaction,
	id uint,
) (*E, error) {
	dbtx := s.GetTransaction(tx).WithContext(ctx)
	var data T
	if err := dbtx.First(&data, id).Error; err != nil {
		return nil, GenerateError("failed to get data", err)
	}

	return s.transformer.ToEntity(&data)
}

func (s *GenericRepository[T, E]) GetMany(
	ctx context.Context,
	tx entities.Transaction,
	ids []uint,
) ([]E, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("%w - input ids is empty", entities.ErrInvalid)
	}

	dbtx := s.GetTransaction(tx).WithContext(ctx)
	var dataArray []T
	if err := dbtx.Find(&dataArray, ids).Error; err != nil {
		return nil, GenerateError("failed to get records", err)
	}

	return s.transformer.ToEntityArray_I2I(dataArray)
}

func (s *GenericRepository[T, E]) GetByCriterias(
	ctx context.Context,
	tx entities.Transaction,
	fields []string,
	criterias map[string]any,
	orderBys []string,
) (*E, error) {
	var data T
	dbtx := s.
		GetTransaction(tx).
		WithContext(ctx)
	if len(fields) > 0 {
		dbtx = dbtx.Select(fields)
	}
	for k, v := range criterias {
		if v == nil {
			dbtx = dbtx.Where(k)
		} else {
			dbtx = dbtx.Where(k, v)
		}
	}
	for _, order := range orderBys {
		dbtx = dbtx.Order(order)
	}

	if err := dbtx.First(&data).Error; err != nil {
		return nil, GenerateError("failed to get data", err)
	}

	return s.transformer.ToEntity(&data)
}

func (s *GenericRepository[T, E]) GetManyByCriterias(
	ctx context.Context,
	tx entities.Transaction,
	fields []string,
	criterias map[string]any,
	orderBys []string,
	offset int,
	limit int,
) ([]E, error) {
	dbtx := s.GetTransaction(tx).WithContext(ctx)

	for k, v := range criterias {
		if v == nil {
			dbtx = dbtx.Where(k)
		} else {
			dbtx = dbtx.Where(k, v)
		}
	}
	for _, order := range orderBys {
		dbtx = dbtx.Order(order)
	}

	if limit <= 0 {
		limit = DefaultLimit
	}
	var dataArray []T
	if err := dbtx.
		Offset(offset).
		Limit(limit).
		Select(fields).
		Find(&dataArray).
		Error; err != nil {
		return nil, GenerateError("failed to get data records", err)
	}

	return s.transformer.ToEntityArray_I2I(dataArray)
}

func (s *GenericRepository[T, E]) Count(
	ctx context.Context,
	tx entities.Transaction,
	criterias map[string]any,
) (int64, error) {
	dbtx := s.GetTransaction(tx).WithContext(ctx)

	for k, v := range criterias {
		if v == nil {
			dbtx = dbtx.Where(k)
		} else {
			dbtx = dbtx.Where(k, v)
		}
	}

	var data T
	var cnt int64
	if err := dbtx.Model(&data).Count(&cnt).Error; err != nil {
		return 0, GenerateError("failed to count", err)
	}

	return cnt, nil
}

func (s *GenericRepository[T, E]) Update(
	ctx context.Context,
	tx entities.Transaction,
	entity *E,
) error {
	if entity == nil {
		return fmt.Errorf("%w - input data is nil", entities.ErrInvalid)
	}

	data, err := s.transformer.FromEntity(entity)
	if err != nil {
		return err
	}

	dbtx := s.GetTransaction(tx).WithContext(ctx)
	dbtx = dbtx.Updates(data)
	if err := dbtx.Error; err != nil {
		return GenerateError("failed to update data", err)
	}
	if dbtx.RowsAffected == 0 {
		return fmt.Errorf("%w - no rows affected", entities.ErrNotFound)
	}

	return nil
}

func (s *GenericRepository[T, E]) Delete(
	ctx context.Context,
	tx entities.Transaction,
	permanent bool,
	id uint,
) error {
	dbtx := s.GetTransaction(tx).WithContext(ctx)
	if permanent {
		dbtx = dbtx.Unscoped()
	}

	var data T
	if err := dbtx.
		Model(&data).
		Delete("id = ?", id).Error; err != nil {
		return GenerateError("failed to delete data", err)
	}

	return nil
}

func (s *GenericRepository[T, E]) DeleteMany(
	ctx context.Context,
	tx entities.Transaction,
	permanent bool,
	ids []uint,
) (int64, error) {
	if len(ids) == 0 {
		return 0, fmt.Errorf("%w - input ids is empty", entities.ErrInvalid)
	}

	dbtx := s.GetTransaction(tx).WithContext(ctx)
	if permanent {
		dbtx = dbtx.Unscoped()
	}

	var data T
	dbtx = dbtx.Model(&data).Delete("id in (?)", ids)
	if err := dbtx.Error; err != nil {
		return 0, GenerateError("failed to delete data", err)
	}

	return dbtx.RowsAffected, nil
}
