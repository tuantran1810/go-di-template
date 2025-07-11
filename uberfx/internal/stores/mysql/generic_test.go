package mysql

import (
	"context"
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type Data struct {
	gorm.Model
	UniqueID string `gorm:"size:32;uniqueIndex"`
	Key      string
	Value    string
}

func setup(t *testing.T) (*GenericStore[Data], error) {
	t.Helper()

	config := RepositoryConfig{
		Username:  "root",
		Password:  "secret",
		Protocol:  "tcp",
		Address:   "127.0.0.1:3306",
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
		return nil, err
	}

	store := NewGenericStore[Data](r)
	if err := store.AutoMigrate(context.Background()); err != nil {
		return nil, err
	}

	return store, nil
}

func cleanup(t *testing.T, store *GenericStore[Data]) {
	t.Helper()

	if err := store.repository.db.Exec("DROP TABLE IF EXISTS `test`.`data`").Error; err != nil {
		t.Logf("failed to cleanup data: %v\n", err)
		return
	}
}

func getTestData(t *testing.T) []Data {
	t.Helper()

	now := time.Now().UTC().Truncate(time.Second)
	return []Data{
		{
			Model: gorm.Model{
				CreatedAt: now,
				UpdatedAt: now,
			},
			UniqueID: "unique-id-1",
			Key:      "key1",
			Value:    "value1",
		},
		{
			Model: gorm.Model{
				CreatedAt: now,
				UpdatedAt: now,
			},
			UniqueID: "unique-id-2",
			Key:      "key2",
			Value:    "value2",
		},
		{
			Model: gorm.Model{
				CreatedAt: now,
				UpdatedAt: now,
			},
			UniqueID: "unique-id-3",
			Key:      "key3",
			Value:    "value3",
		},
	}
}

func createTestData(t *testing.T, store *GenericStore[Data]) ([]Data, error) {
	t.Helper()
	testData := getTestData(t)
	_, err := store.CreateMany(context.Background(), nil, testData)
	return testData, err
}

type GenericDataTestSuite struct {
	suite.Suite
	initData []Data
	store    *GenericStore[Data]
}

func (s *GenericDataTestSuite) SetupSuite() {
	t := s.T()
	os.Setenv("TZ", "UTC")
	store, err := setup(t)
	s.Require().NoError(err)
	s.store = store
	s.Require().NotNil(s.store)
	s.Require().NotNil(s.store.repository)
	s.Require().NotNil(s.store.repository.db)
	if err := store.repository.db.Exec("SET @@global.time_zone = '+00:00'").Error; err != nil {
		t.Errorf("failed to set time zone: %v", err)
		return
	}
}

func (s *GenericDataTestSuite) TearDownSuite() {
	t := s.T()
	cleanup(t, s.store)
}

func (s *GenericDataTestSuite) SetupTest() {
	t := s.T()
	if err := s.store.repository.db.Exec("TRUNCATE TABLE `test`.`data`").Error; err != nil {
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
	if err := s.store.repository.db.Exec("TRUNCATE TABLE `test`.`data`").Error; err != nil {
		t.Errorf("failed to cleanup data: %v\n", err)
		return
	}
}

func (s *GenericDataTestSuite) TestGenericStore_PingOK() {
	t := s.T()

	t.Run("Ping", func(t *testing.T) {
		if err := s.store.Ping(context.Background()); err != nil {
			t.Errorf("failed to ping database: %v", err)
		}
	})
}

func (s *GenericDataTestSuite) TestGenericStore_PingFailed() {
	t := s.T()
	s.store.repository.db.Exec("DROP TABLE IF EXISTS `test`.`data`")
	t.Run("Ping", func(t *testing.T) {
		if err := s.store.Ping(context.Background()); err == nil {
			t.Error("expected error")
		}
	})
}

func (s *GenericDataTestSuite) TestGenericStore_Create() {
	t := s.T()
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name    string
		input   *Data
		want    *Data
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
			input: &Data{
				Model: gorm.Model{
					CreatedAt: now,
					UpdatedAt: now,
				},
				UniqueID: "unique-id-4",
				Key:      "key",
				Value:    "value",
			},
			want: &Data{
				Model: gorm.Model{
					ID:        4,
					CreatedAt: now,
					UpdatedAt: now,
				},
				UniqueID: "unique-id-4",
				Key:      "key",
				Value:    "value",
			},
			wantErr: false,
		},
		{
			name: "no error, key 5",
			input: &Data{
				Model: gorm.Model{
					CreatedAt: now,
					UpdatedAt: now,
				},
				UniqueID: "unique-id-5",
				Key:      "key",
				Value:    "value",
			},
			want: &Data{
				Model: gorm.Model{
					ID:        5,
					CreatedAt: now,
					UpdatedAt: now,
				},
				UniqueID: "unique-id-5",
				Key:      "key",
				Value:    "value",
			},
			wantErr: false,
		},
		{
			name: "error, conflicted",
			input: &Data{
				Model: gorm.Model{
					CreatedAt: now,
					UpdatedAt: now,
				},
				UniqueID: "unique-id-5",
				Key:      "key",
				Value:    "value",
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

func (s *GenericDataTestSuite) TestGenericStore_CreateMany() {
	t := s.T()
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name    string
		input   []Data
		want    []Data
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
			input: []Data{
				{
					Model: gorm.Model{
						CreatedAt: now,
						UpdatedAt: now,
					},
					UniqueID: "unique-id-4",
					Key:      "key4",
					Value:    "value4",
				},
				{
					Model: gorm.Model{
						CreatedAt: now,
						UpdatedAt: now,
					},
					UniqueID: "unique-id-5",
					Key:      "key5",
					Value:    "value5",
				},
			},
			want: []Data{
				{
					Model: gorm.Model{
						ID:        4,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UniqueID: "unique-id-4",
					Key:      "key4",
					Value:    "value4",
				},
				{
					Model: gorm.Model{
						ID:        5,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UniqueID: "unique-id-5",
					Key:      "key5",
					Value:    "value5",
				},
			},
			wantErr: false,
		},
		{
			name: "error, conflicted",
			input: []Data{
				{
					Model: gorm.Model{
						CreatedAt: now,
						UpdatedAt: now,
					},
					UniqueID: "unique-id-2",
					Key:      "key2",
					Value:    "value2",
				},
				{
					Model: gorm.Model{
						CreatedAt: now,
						UpdatedAt: now,
					},
					UniqueID: "unique-id-3",
					Key:      "key3",
					Value:    "value3",
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

func (s *GenericDataTestSuite) TestGenericStore_Get() {
	t := s.T()

	tests := []struct {
		name    string
		id      uint
		want    *Data
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

func (s *GenericDataTestSuite) TestGenericStore_GetMany() {
	t := s.T()

	tests := []struct {
		name    string
		ids     []uint
		want    []*Data
		wantErr bool
	}{
		{
			name:    "key 1,2",
			ids:     []uint{1, 2},
			want:    []*Data{&s.initData[0], &s.initData[1]},
			wantErr: false,
		},
		{
			name:    "key 1,3",
			ids:     []uint{1, 3},
			want:    []*Data{&s.initData[0], &s.initData[2]},
			wantErr: false,
		},
		{
			name:    "key 1,10",
			ids:     []uint{1, 10},
			want:    []*Data{&s.initData[0]},
			wantErr: false,
		},
		{
			name:    "not found",
			ids:     []uint{10, 11},
			want:    []*Data{},
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

func (s *GenericDataTestSuite) TestGenericStore_GetByCriterias() {
	t := s.T()

	tests := []struct {
		name      string
		criterias map[string]any
		fields    []string
		orderBys  []string
		want      *Data
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
			want: &Data{
				UniqueID: "unique-id-1",
				Key:      "key1",
			},
			wantErr: false,
		},
		{
			name:     "without criterias, limited fields, order by",
			fields:   []string{"unique_id", "key"},
			orderBys: []string{"id DESC"},
			want: &Data{
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
			want: &Data{
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

func (s *GenericDataTestSuite) TestGenericStore_GetManyByCriterias() {
	t := s.T()
	now := time.Now().UTC().Truncate(time.Second)

	if _, err := s.store.Create(context.Background(), nil, &Data{
		Model: gorm.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
		UniqueID: "unique-id-4",
		Key:      "key1",
		Value:    "value2",
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
		want      []*Data
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
			want: []*Data{
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
			want: []*Data{
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
			want: []*Data{
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
			want:    []*Data{},
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

func (s *GenericDataTestSuite) TestGenericStore_Update() {
	t := s.T()
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name    string
		data    *Data
		wantErr bool
	}{
		{
			name: "key 1",
			data: &Data{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				UniqueID: "unique-id-1",
				Key:      "key1_updated",
				Value:    "value1_updated",
			},
			wantErr: false,
		},
		{
			name: "not found",
			data: &Data{
				Model: gorm.Model{
					ID:        100,
					CreatedAt: now,
					UpdatedAt: now,
				},
				UniqueID: "unique-id-100",
				Key:      "key100_updated",
				Value:    "value100_updated",
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

func (s *GenericDataTestSuite) TestGenericStore_Delete() {
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
			} else if err := s.store.repository.db.Unscoped().First(&Data{}, tt.id).Error; err == nil {
				t.Errorf("still found after delete")
				return
			}
		})
	}
}

func (s *GenericDataTestSuite) TestGenericStore_DeleteMany() {
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
				if err := s.store.repository.db.Unscoped().Find(&data, tt.ids).Error; err == nil && len(data) != 0 {
					t.Errorf("still found after delete")
					return
				}
			}
		})
	}
}

func TestGenericDataTestSuite(t *testing.T) {
	suite.Run(t, new(GenericDataTestSuite))
}
