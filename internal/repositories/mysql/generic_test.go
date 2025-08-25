package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tuantran1810/go-di-template/internal/entities"
	"gorm.io/gorm"
)

type Data struct {
	gorm.Model
	UniqueID string `gorm:"size:32;uniqueIndex"`
	Key      string
	Value    string
}

type DataEntity struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	UniqueID  string
	Key       string
	Value     string
}

type DataTransformer struct{}

func (t *DataTransformer) ToEntity(data *Data) (*DataEntity, error) {
	if data == nil {
		return nil, fmt.Errorf("%w - input data is nil", entities.ErrInvalid)
	}

	return &DataEntity{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		UniqueID:  data.UniqueID,
		Key:       data.Key,
		Value:     data.Value,
	}, nil
}

func (t *DataTransformer) FromEntity(entity *DataEntity) (*Data, error) {
	if entity == nil {
		return nil, fmt.Errorf("%w - input entity is nil", entities.ErrInvalid)
	}

	return &Data{
		Model: gorm.Model{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		UniqueID: entity.UniqueID,
		Key:      entity.Key,
		Value:    entity.Value,
	}, nil
}

type FkData struct {
	gorm.Model
	DataRefer uint
	Data      Data `gorm:"foreignKey:DataRefer"`
	Metadata  string
}

type FkDataEntity struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DataRefer uint
	Data      DataEntity
	Metadata  string
}

type FkDataTransformer struct{}

func (t *FkDataTransformer) ToEntity(data *FkData) (*FkDataEntity, error) {
	if data == nil {
		return nil, fmt.Errorf("%w - input data is nil", entities.ErrInvalid)
	}

	var dataTransformer DataTransformer
	tmp, err := dataTransformer.ToEntity(&data.Data)
	if err != nil {
		return nil, err
	}

	return &FkDataEntity{
		ID:        data.ID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		DataRefer: data.DataRefer,
		Data:      *tmp,
		Metadata:  data.Metadata,
	}, nil
}

func (t *FkDataTransformer) FromEntity(entity *FkDataEntity) (*FkData, error) {
	if entity == nil {
		return nil, fmt.Errorf("%w - input entity is nil", entities.ErrInvalid)
	}

	return &FkData{
		Model: gorm.Model{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		DataRefer: entity.DataRefer,
		Metadata:  entity.Metadata,
	}, nil
}

type DataStore = GenericRepository[Data, DataEntity]
type FkDataStore = GenericRepository[FkData, FkDataEntity]

func setup(t *testing.T, port int) (*DataStore, *FkDataStore, error) {
	t.Helper()

	config := RepositoryConfig{
		Username:  "root",
		Password:  "secret",
		Protocol:  "tcp",
		Address:   fmt.Sprintf("127.0.0.1:%d", port),
		Database:  "test",
		Params:    map[string]string{},
		Collation: "utf8mb4_general_ci",
		Loc:       time.Local,
		TLSConfig: "",

		Timeout:                 10 * time.Second,
		ReadTimeout:             10 * time.Second,
		WriteTimeout:            10 * time.Second,
		AllowAllFiles:           false,
		AllowCleartextPasswords: false,
		AllowOldPasswords:       false,
		ClientFoundRows:         false,
		ColumnsWithAlias:        false,
		InterpolateParams:       false,
		MultiStatements:         false,
		ParseTime:               true,

		MaxOpenConns:           10,
		MaxIdleConns:           10,
		ConnMaxLifeTimeSeconds: 1800,
	}
	r := MustNewRepository(config)
	if err := r.Start(context.Background()); err != nil {
		return nil, nil, err
	}

	dataTransformer := entities.NewExtendedDataTransformer(&DataTransformer{})
	store := NewGenericRepository(r, dataTransformer)
	if err := store.AutoMigrate(context.Background()); err != nil {
		return nil, nil, err
	}

	fkDataTransformer := entities.NewExtendedDataTransformer(&FkDataTransformer{})
	fkStore := NewGenericRepository(r, fkDataTransformer)
	if err := fkStore.AutoMigrate(context.Background()); err != nil {
		return nil, nil, err
	}

	return store, fkStore, nil
}

func cleanup(t *testing.T, store *GenericRepository[Data, DataEntity]) {
	t.Helper()

	if err := store.db.Exec("DROP TABLE IF EXISTS `test`.`fk_data`").Error; err != nil {
		t.Logf("failed to cleanup fk_data: %v\n", err)
		return
	}

	if err := store.db.Exec("DROP TABLE IF EXISTS `test`.`data`").Error; err != nil {
		t.Logf("failed to cleanup data: %v\n", err)
		return
	}
}

func getTestData(t *testing.T) []DataEntity {
	t.Helper()

	now := time.Now().UTC().Truncate(time.Second)
	return []DataEntity{
		{
			CreatedAt: now,
			UpdatedAt: now,
			UniqueID:  "unique-id-1",
			Key:       "key1",
			Value:     "value1",
		},
		{
			CreatedAt: now,
			UpdatedAt: now,
			UniqueID:  "unique-id-2",
			Key:       "key2",
			Value:     "value2",
		},
		{
			CreatedAt: now,
			UpdatedAt: now,
			UniqueID:  "unique-id-3",
			Key:       "key3",
			Value:     "value3",
		},
	}
}

func createTestData(t *testing.T, store *DataStore) ([]DataEntity, error) {
	t.Helper()
	testData := getTestData(t)
	return store.CreateMany(context.Background(), nil, testData)
}

type GenericDataTestSuite struct {
	suite.Suite
	initData  []DataEntity
	container *mysql.MySQLContainer
	store     *DataStore
	fkStore   *FkDataStore
}

func (s *GenericDataTestSuite) SetupSuite() {
	t := s.T()
	t.Helper()
	if err := os.Setenv("TZ", "UTC"); err != nil {
		t.Errorf("failed to set time zone: %v", err)
		return
	}

	mysqlContainer, err := mysql.Run(context.Background(),
		"mysql:lts",
		mysql.WithDatabase("test"),
		mysql.WithUsername("root"),
		mysql.WithPassword("secret"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("port: 3306  MySQL Community Server - GPL").WithStartupTimeout(30*time.Second),
			wait.ForListeningPort("3306/tcp").WithStartupTimeout(30*time.Second),
		),
	)
	s.Require().NoError(err)

	port, err := mysqlContainer.MappedPort(context.Background(), "3306")
	s.Require().NoError(err)
	s.Require().NotNil(port)

	s.Require().NoError(err)
	s.container = mysqlContainer
	s.Require().NotNil(s.container)

	store, fkStore, err := setup(t, port.Int())
	s.Require().NoError(err)
	s.store = store
	s.Require().NotNil(s.store)
	s.Require().NotNil(s.store.Repository)
	s.Require().NotNil(s.store.db)
	if err := store.db.Exec("SET @@global.time_zone = '+00:00'").Error; err != nil {
		t.Errorf("failed to set time zone: %v", err)
		return
	}

	s.fkStore = fkStore
	s.Require().NotNil(s.fkStore)
	s.Require().NotNil(s.fkStore.Repository)
	s.Require().NotNil(s.fkStore.db)
}

func (s *GenericDataTestSuite) TearDownSuite() {
	t := s.T()
	t.Helper()
	cleanup(t, s.store)

	if err := testcontainers.TerminateContainer(s.container); err != nil {
		t.Errorf("failed to terminate container: %v", err)
		return
	}
}

func (s *GenericDataTestSuite) SetupTest() {
	t := s.T()
	t.Helper()
	s.Require().NoError(s.store.AutoMigrate(context.Background()))

	if err := s.store.db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		t.Errorf("failed to cleanup data: %v\n", err)
		return
	}

	if err := s.fkStore.db.Exec("TRUNCATE TABLE `test`.`fk_data`").Error; err != nil {
		t.Errorf("failed to cleanup fk_data: %v\n", err)
		return
	}

	if err := s.store.db.Exec("TRUNCATE TABLE `test`.`data`").Error; err != nil {
		t.Errorf("failed to cleanup data: %v\n", err)
		return
	}

	if err := s.store.db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		t.Errorf("failed to cleanup data: %v\n", err)
		return
	}

	data, err := createTestData(t, s.store)
	if err != nil {
		t.Errorf("failed to create test data: %v\n", err)
		return
	}

	s.initData = data
}

func (s *GenericDataTestSuite) TearDownTest() {
	t := s.T()
	t.Helper()

	if err := s.store.db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		t.Errorf("failed to cleanup data: %v\n", err)
		return
	}

	if err := s.fkStore.db.Exec("TRUNCATE TABLE `test`.`fk_data`").Error; err != nil {
		t.Errorf("failed to cleanup fk_data: %v\n", err)
		return
	}

	if err := s.store.db.Exec("TRUNCATE TABLE `test`.`data`").Error; err != nil {
		t.Errorf("failed to cleanup data: %v\n", err)
		return
	}

	if err := s.store.db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		t.Errorf("failed to cleanup data: %v\n", err)
		return
	}
}

func (s *GenericDataTestSuite) TestGenericRepository_PingOK() {
	t := s.T()

	t.Run("Ping", func(t *testing.T) {
		if err := s.store.Ping(context.Background()); err != nil {
			t.Errorf("failed to ping database: %v", err)
		}
	})
}

func (s *GenericDataTestSuite) TestGenericRepository_PingFailed() {
	t := s.T()
	s.fkStore.db.Exec("DROP TABLE IF EXISTS `test`.`fk_data`")
	t.Run("Ping", func(t *testing.T) {
		if err := s.fkStore.Ping(context.Background()); err == nil {
			t.Error("expected error")
		}
	})

	s.Require().NoError(s.fkStore.AutoMigrate(context.Background()))
}

func (s *GenericDataTestSuite) TestGenericRepository_Create() {
	t := s.T()
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name    string
		input   *DataEntity
		want    *DataEntity
		wantErr bool
	}{
		{
			name:    "nil input",
			input:   nil,
			want:    nil,
			wantErr: true,
		},
		{
			name: "no error, key 4",
			input: &DataEntity{
				CreatedAt: now,
				UpdatedAt: now,
				UniqueID:  "unique-id-4",
				Key:       "key",
				Value:     "value",
			},
			want: &DataEntity{
				ID:        4,
				CreatedAt: now,
				UpdatedAt: now,
				UniqueID:  "unique-id-4",
				Key:       "key",
				Value:     "value",
			},
			wantErr: false,
		},
		{
			name: "no error, key 5",
			input: &DataEntity{
				CreatedAt: now,
				UpdatedAt: now,
				UniqueID:  "unique-id-5",
				Key:       "key",
				Value:     "value",
			},
			want: &DataEntity{
				ID:        5,
				CreatedAt: now,
				UpdatedAt: now,
				UniqueID:  "unique-id-5",
				Key:       "key",
				Value:     "value",
			},
			wantErr: false,
		},
		{
			name: "error, conflicted",
			input: &DataEntity{
				CreatedAt: now,
				UpdatedAt: now,
				UniqueID:  "unique-id-5",
				Key:       "key",
				Value:     "value",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.store.Create(context.Background(), nil, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("store.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *GenericDataTestSuite) TestGenericRepository_CreateConflict() {
	t := s.T()
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name    string
		input   *FkDataEntity
		want    *FkDataEntity
		wantErr bool
	}{
		{
			name: "error, no fk data",
			input: &FkDataEntity{
				CreatedAt: now,
				UpdatedAt: now,
				DataRefer: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.fkStore.Create(context.Background(), nil, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("fkStore.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fkStore.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *GenericDataTestSuite) TestGenericRepository_CreateMany() {
	t := s.T()
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name    string
		input   []DataEntity
		want    []DataEntity
		wantErr bool
	}{
		{
			name:    "nil input",
			input:   nil,
			want:    nil,
			wantErr: true,
		},
		{
			name: "no error, key 4-5",
			input: []DataEntity{
				{
					CreatedAt: now,
					UpdatedAt: now,
					UniqueID:  "unique-id-4",
					Key:       "key4",
					Value:     "value4",
				},
				{
					CreatedAt: now,
					UpdatedAt: now,
					UniqueID:  "unique-id-5",
					Key:       "key5",
					Value:     "value5",
				},
			},
			want: []DataEntity{
				{
					ID:        4,
					CreatedAt: now,
					UpdatedAt: now,
					UniqueID:  "unique-id-4",
					Key:       "key4",
					Value:     "value4",
				},
				{
					ID:        5,
					CreatedAt: now,
					UpdatedAt: now,
					UniqueID:  "unique-id-5",
					Key:       "key5",
					Value:     "value5",
				},
			},
			wantErr: false,
		},
		{
			name: "error, conflicted",
			input: []DataEntity{
				{
					CreatedAt: now,
					UpdatedAt: now,
					UniqueID:  "unique-id-2",
					Key:       "key2",
					Value:     "value2",
				},
				{
					CreatedAt: now,
					UpdatedAt: now,
					UniqueID:  "unique-id-3",
					Key:       "key3",
					Value:     "value3",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.store.CreateMany(context.Background(), nil, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.CreateMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("store.CreateMany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *GenericDataTestSuite) TestGenericRepository_Get() {
	t := s.T()

	tests := []struct {
		name    string
		id      uint
		want    *DataEntity
		wantErr bool
	}{
		{
			name:    "key 1",
			id:      1,
			want:    &s.initData[0],
			wantErr: false,
		},
		{
			name:    "key 2",
			id:      2,
			want:    &s.initData[1],
			wantErr: false,
		},
		{
			name:    "not found",
			id:      10,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.store.Get(context.Background(), nil, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("store.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *GenericDataTestSuite) TestGenericRepository_GetMany() {
	t := s.T()

	tests := []struct {
		name    string
		ids     []uint
		want    []DataEntity
		wantErr bool
	}{
		{
			name:    "key 1,2",
			ids:     []uint{1, 2},
			want:    []DataEntity{s.initData[0], s.initData[1]},
			wantErr: false,
		},
		{
			name:    "key 1,3",
			ids:     []uint{1, 3},
			want:    []DataEntity{s.initData[0], s.initData[2]},
			wantErr: false,
		},
		{
			name:    "key 1,10",
			ids:     []uint{1, 10},
			want:    []DataEntity{s.initData[0]},
			wantErr: false,
		},
		{
			name:    "not found",
			ids:     []uint{10, 11},
			want:    []DataEntity{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.store.GetMany(context.Background(), nil, tt.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.GetMany() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("store.GetMany() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *GenericDataTestSuite) TestGenericRepository_GetByCriterias() {
	t := s.T()

	tests := []struct {
		name      string
		criterias map[string]any
		fields    []string
		orderBys  []string
		want      *DataEntity
		wantErr   bool
	}{
		{
			name: "with criterias",
			criterias: map[string]any{
				"unique_id": "unique-id-1",
			},
			orderBys: []string{"id"},
			want:     &s.initData[0],
			wantErr:  false,
		},
		{
			name: "with criterias, limited fields",
			criterias: map[string]any{
				"unique_id": "unique-id-1",
			},
			fields:   []string{"unique_id", "key"},
			orderBys: []string{"id"},
			want: &DataEntity{
				UniqueID: "unique-id-1",
				Key:      "key1",
			},
			wantErr: false,
		},
		{
			name:     "without criterias, limited fields, order by",
			fields:   []string{"unique_id", "key"},
			orderBys: []string{"id DESC"},
			want: &DataEntity{
				UniqueID: "unique-id-3",
				Key:      "key3",
			},
			wantErr: false,
		},
		{
			name: "not found",
			criterias: map[string]any{
				"unique_id": "unique-id-100",
			},
			fields:   []string{"unique_id", "key"},
			orderBys: []string{"id"},
			want:     nil,
			wantErr:  true,
		},
		{
			name: "no criterias",
			fields: []string{
				"unique_id",
				"key",
			},
			orderBys: []string{
				"id",
			},
			want: &DataEntity{
				UniqueID: "unique-id-1",
				Key:      "key1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.store.GetByCriterias(context.Background(), nil, tt.fields, tt.criterias, tt.orderBys)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.GetByCriterias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("store.GetByCriterias() = %v, want %v", got, tt.want)
			}
		})
	}
}

func (s *GenericDataTestSuite) TestGenericRepository_GetManyByCriterias() {
	t := s.T()
	now := time.Now().UTC().Truncate(time.Second)

	if _, err := s.store.Create(context.Background(), nil, &DataEntity{
		CreatedAt: now,
		UpdatedAt: now,
		UniqueID:  "unique-id-4",
		Key:       "key1",
		Value:     "value2",
	}); err != nil {
		t.Errorf("failed to create data: %v", err)
		return
	}

	tests := []struct {
		name      string
		criterias map[string]any
		fields    []string
		orderBys  []string
		offset    int
		limit     int
		want      []DataEntity
		wantErr   bool
	}{
		{
			name: "key 1",
			criterias: map[string]any{
				"key": "key1",
			},
			fields: []string{"unique_id", "key"},
			orderBys: []string{
				"id DESC",
			},
			offset: 0,
			limit:  10,
			want: []DataEntity{
				{
					UniqueID: "unique-id-4",
					Key:      "key1",
				},
				{
					UniqueID: "unique-id-1",
					Key:      "key1",
				},
			},
			wantErr: false,
		},
		{
			name: "key 1, offset 1",
			criterias: map[string]any{
				"key": "key1",
			},
			fields: []string{"unique_id", "key"},
			orderBys: []string{
				"id DESC",
			},
			offset: 1,
			limit:  10,
			want: []DataEntity{
				{
					UniqueID: "unique-id-1",
					Key:      "key1",
				},
			},
			wantErr: false,
		},
		{
			name: "key 1, limit 1",
			criterias: map[string]any{
				"key": "key1",
			},
			fields: []string{"unique_id", "key"},
			orderBys: []string{
				"id DESC",
			},
			offset: 0,
			limit:  1,
			want: []DataEntity{
				{
					UniqueID: "unique-id-4",
					Key:      "key1",
				},
			},
			wantErr: false,
		},
		{
			name: "key 1, offset 2",
			criterias: map[string]any{
				"key": "key1",
			},
			fields: []string{"unique_id", "key"},
			orderBys: []string{
				"id DESC",
			},
			offset:  2,
			limit:   10,
			want:    []DataEntity{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.store.GetManyByCriterias(
				context.Background(), nil,
				tt.fields, tt.criterias, tt.orderBys,
				tt.offset, tt.limit,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.GetManyByCriterias() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				jgot, _ := json.Marshal(got)
				jwant, _ := json.Marshal(tt.want)
				t.Errorf("store.GetManyByCriterias() = %s, want %s", jgot, jwant)
			}
		})
	}
}

func (s *GenericDataTestSuite) TestGenericRepository_Count() {
	t := s.T()
	now := time.Now().UTC().Truncate(time.Second)

	if _, err := s.store.Create(context.Background(), nil, &DataEntity{
		CreatedAt: now,
		UpdatedAt: now,
		UniqueID:  "unique-id-4",
		Key:       "key1",
		Value:     "value2",
	}); err != nil {
		t.Errorf("failed to create data: %v", err)
		return
	}

	tests := []struct {
		name      string
		criterias map[string]any
		want      int64
		wantErr   bool
	}{
		{
			name: "key 1",
			criterias: map[string]any{
				"key": "key1",
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "key 2",
			criterias: map[string]any{
				"key": "key2",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "key 1000",
			criterias: map[string]any{
				"key": "key1000",
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.store.Count(
				context.Background(), nil,
				tt.criterias,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				jgot, _ := json.Marshal(got)
				jwant, _ := json.Marshal(tt.want)
				t.Errorf("store.Count() = %s, want %s", jgot, jwant)
			}
		})
	}
}

func (s *GenericDataTestSuite) TestGenericRepository_Update() {
	t := s.T()
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name    string
		data    *DataEntity
		wantErr bool
	}{
		{
			name: "key 1",
			data: &DataEntity{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
				UniqueID:  "unique-id-1",
				Key:       "key1_updated",
				Value:     "value1_updated",
			},
			wantErr: false,
		},
		{
			name: "not found",
			data: &DataEntity{
				ID:        100,
				CreatedAt: now,
				UpdatedAt: now,
				UniqueID:  "unique-id-100",
				Key:       "key100_updated",
				Value:     "value100_updated",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.store.Update(context.Background(), nil, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func (s *GenericDataTestSuite) TestGenericRepository_Delete() {
	t := s.T()

	tests := []struct {
		name      string
		id        uint
		permanent bool
		wantErr   bool
	}{
		{
			name:      "key 1",
			id:        1,
			permanent: false,
			wantErr:   false,
		},
		{
			name:      "key 2",
			id:        2,
			permanent: false,
		},
		{
			name:      "key 3, permanent",
			id:        3,
			permanent: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.store.Delete(context.Background(), nil, tt.permanent, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.permanent {
				if _, err := s.store.Get(context.Background(), nil, tt.id); err == nil {
					t.Errorf("still found after delete")
					return
				}
			} else if err := s.store.db.Unscoped().First(&Data{}, tt.id).Error; err == nil {
				t.Errorf("still found after delete")
				return
			}
		})
	}
}

func (s *GenericDataTestSuite) TestGenericRepository_DeleteMany() {
	t := s.T()

	tests := []struct {
		name      string
		ids       []uint
		permanent bool
		want      int64
		wantErr   bool
	}{
		{
			name:      "key 1, 2",
			ids:       []uint{1, 2},
			permanent: false,
			want:      2,
			wantErr:   false,
		},
		{
			name:      "key 1, 2, permanent",
			ids:       []uint{1, 2},
			permanent: true,
			want:      2,
			wantErr:   false,
		},
		{
			name:      "key 2,3 permanent",
			ids:       []uint{2, 3},
			permanent: true,
			want:      1,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.store.DeleteMany(context.Background(), nil, tt.permanent, tt.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("store.DeleteMany() = %v, want %v", got, tt.want)
			}

			if !tt.permanent {
				if out, err := s.store.GetMany(context.Background(), nil, tt.ids); err == nil && len(out) != 0 {
					t.Errorf("still found after delete")
					return
				}
			} else {
				var data []*Data
				if err := s.store.db.Unscoped().Find(&data, tt.ids).Error; err == nil && len(data) != 0 {
					t.Errorf("still found after delete")
					return
				}
			}
		})
	}
}

func (s *GenericDataTestSuite) TestGenericRepository_Transaction() {
	t := s.T()

	createUniqueId4 := func(ctx context.Context, txKeeper entities.Transaction) (*DataEntity, error) {
		return s.store.Create(ctx, txKeeper, &DataEntity{
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UpdatedAt: time.Now().UTC().Truncate(time.Second),
			UniqueID:  "unique-id-4",
			Key:       "key4",
			Value:     "value4",
		})
	}

	createUniqueId5 := func(ctx context.Context, txKeeper entities.Transaction) (*DataEntity, error) {
		return s.store.Create(ctx, txKeeper, &DataEntity{
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UpdatedAt: time.Now().UTC().Truncate(time.Second),
			UniqueID:  "unique-id-5",
			Key:       "key5",
			Value:     "value5",
		})
	}

	createUniqueId6 := func(ctx context.Context, txKeeper entities.Transaction) (*DataEntity, error) {
		return s.store.Create(ctx, txKeeper, &DataEntity{
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UpdatedAt: time.Now().UTC().Truncate(time.Second),
			UniqueID:  "unique-id-6",
			Key:       "key6",
			Value:     "value6",
		})
	}

	createUniqueId7 := func(ctx context.Context, txKeeper entities.Transaction) (*DataEntity, error) {
		return s.store.Create(ctx, txKeeper, &DataEntity{
			CreatedAt: time.Now().UTC().Truncate(time.Second),
			UpdatedAt: time.Now().UTC().Truncate(time.Second),
			UniqueID:  "unique-id-7",
			Key:       "key7",
			Value:     "value7",
		})
	}

	tests := []struct {
		name      string
		funcs     []entities.DBTxHandleFunc
		want      any
		wantCount int
		wantErr   bool
	}{
		{
			name: "no error",
			funcs: []entities.DBTxHandleFunc{
				func(ctx context.Context, txKeeper entities.Transaction) error {
					if _, err := createUniqueId4(ctx, txKeeper); err != nil {
						return err
					}
					if _, err := createUniqueId5(ctx, txKeeper); err != nil {
						return err
					}
					return nil
				},
			},
			wantCount: 5,
			wantErr:   false,
		},
		{
			name: "with error",
			funcs: []entities.DBTxHandleFunc{
				func(ctx context.Context, txKeeper entities.Transaction) error {
					if _, err := createUniqueId6(ctx, txKeeper); err != nil {
						return err
					}
					if _, err := createUniqueId7(ctx, txKeeper); err != nil {
						return err
					}
					return nil
				},
				func(ctx context.Context, txKeeper entities.Transaction) error {
					return fmt.Errorf("fake error")
				},
			},
			want:      nil,
			wantCount: 5,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.store.RunTx(
				context.Background(),
				tt.funcs...,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.RunTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var cnt int64
			if err := s.store.DB().Model(&Data{}).Count(&cnt).Error; err != nil {
				t.Errorf("failed to count data: %v", err)
				return
			}
			if cnt != int64(tt.wantCount) {
				t.Errorf("store.RunTx() count = %v, want %v", cnt, tt.wantCount)
			}
		})
	}
}

func TestGenericDataTestSuite(t *testing.T) {
	suite.Run(t, new(GenericDataTestSuite))
}
