package sqlite

import (
	"context"
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"
)

type Data struct {
	gorm.Model
	UniqueID string `gorm:"uniqueIndex"`
	Key      string
	Value    string
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

func createTestData(t *testing.T, store *GenericStore[Data]) {
	t.Helper()

	if _, err := store.CreateMany(context.Background(), nil, getTestData(t)); err != nil {
		t.Errorf("failed to create data: %v", err)
		return
	}
}

func setup(t *testing.T) (*Repository, error) {
	t.Helper()
	r, err := NewRepository(RepositoryConfig{DatabasePath: ":memory:"})
	if err != nil {
		return nil, err
	}

	if err := r.db.AutoMigrate(&Data{}); err != nil {
		return nil, err
	}
	return r, nil
}

func TestGenericStore_Create(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	repository, err := setup(t)
	if err != nil {
		t.Errorf("failed to setup repository: %v", err)
		return
	}
	defer repository.Stop(context.Background())

	store := NewGenericStore[Data](repository)

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
			name: "no error, key 1",
			input: &Data{
				Model: gorm.Model{
					CreatedAt: now,
					UpdatedAt: now,
				},
				UniqueID: "unique-id",
				Key:      "key",
				Value:    "value",
			},
			want: &Data{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				UniqueID: "unique-id",
				Key:      "key",
				Value:    "value",
			},
			wantErr: false,
		},
		{
			name: "no error, key 2",
			input: &Data{
				Model: gorm.Model{
					CreatedAt: now,
					UpdatedAt: now,
				},
				UniqueID: "unique-id-2",
				Key:      "key",
				Value:    "value",
			},
			want: &Data{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: now,
					UpdatedAt: now,
				},
				UniqueID: "unique-id-2",
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
				UniqueID: "unique-id",
				Key:      "key",
				Value:    "value",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.Create(context.Background(), nil, tt.input)
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

func TestGenericStore_CreateMany(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	repository, err := setup(t)
	if err != nil {
		t.Errorf("failed to setup repository: %v", err)
		return
	}
	defer repository.Stop(context.Background())

	store := NewGenericStore[Data](repository)

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
			name: "no error, key 1-2",
			input: []Data{
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
			},
			want: []Data{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UniqueID: "unique-id-1",
					Key:      "key1",
					Value:    "value1",
				},
				{
					Model: gorm.Model{
						ID:        2,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UniqueID: "unique-id-2",
					Key:      "key2",
					Value:    "value2",
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
			got, err := store.CreateMany(context.Background(), nil, tt.input)
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

func TestGenericStore_Get(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	repository, err := setup(t)
	if err != nil {
		t.Errorf("failed to setup repository: %v", err)
		return
	}
	defer repository.Stop(context.Background())

	store := NewGenericStore[Data](repository)
	createTestData(t, store)

	tests := []struct {
		name    string
		id      uint
		want    *Data
		wantErr bool
	}{
		{
			name: "key 1",
			id:   1,
			want: &Data{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				UniqueID: "unique-id-1",
				Key:      "key1",
				Value:    "value1",
			},
			wantErr: false,
		},
		{
			name: "key 2",
			id:   2,
			want: &Data{
				Model: gorm.Model{
					ID:        2,
					CreatedAt: now,
					UpdatedAt: now,
				},
				UniqueID: "unique-id-2",
				Key:      "key2",
				Value:    "value2",
			},
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
			got, err := store.Get(context.Background(), nil, tt.id)
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

func TestGenericStore_GetMany(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	repository, err := setup(t)
	if err != nil {
		t.Errorf("failed to setup repository: %v", err)
		return
	}
	defer repository.Stop(context.Background())

	store := NewGenericStore[Data](repository)
	createTestData(t, store)

	tests := []struct {
		name    string
		ids     []uint
		want    []*Data
		wantErr bool
	}{
		{
			name: "key 1,2",
			ids:  []uint{1, 2},
			want: []*Data{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UniqueID: "unique-id-1",
					Key:      "key1",
					Value:    "value1",
				},
				{
					Model: gorm.Model{
						ID:        2,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UniqueID: "unique-id-2",
					Key:      "key2",
					Value:    "value2",
				},
			},
			wantErr: false,
		},
		{
			name: "key 1,3",
			ids:  []uint{1, 3},
			want: []*Data{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UniqueID: "unique-id-1",
					Key:      "key1",
					Value:    "value1",
				},
				{
					Model: gorm.Model{
						ID:        3,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UniqueID: "unique-id-3",
					Key:      "key3",
					Value:    "value3",
				},
			},
			wantErr: false,
		},
		{
			name: "key 1,10",
			ids:  []uint{1, 10},
			want: []*Data{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: now,
						UpdatedAt: now,
					},
					UniqueID: "unique-id-1",
					Key:      "key1",
					Value:    "value1",
				},
			},
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
			got, err := store.GetMany(context.Background(), nil, tt.ids)
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

func TestGenericStore_Update(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	repository, err := setup(t)
	if err != nil {
		t.Errorf("failed to setup repository: %v", err)
		return
	}
	defer repository.Stop(context.Background())

	store := NewGenericStore[Data](repository)
	createTestData(t, store)

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
			err := store.Update(context.Background(), nil, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGenericStore_Delete(t *testing.T) {
	repository, err := setup(t)
	if err != nil {
		t.Errorf("failed to setup repository: %v", err)
		return
	}
	defer repository.Stop(context.Background())

	store := NewGenericStore[Data](repository)
	createTestData(t, store)

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
			err := store.Delete(context.Background(), nil, tt.permanent, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.permanent {
				store.repository.mutex.RLock()
				if _, err := store.Get(context.Background(), nil, tt.id); err == nil {
					t.Errorf("still found after delete")
					store.repository.mutex.RUnlock()
					return
				} else {
					store.repository.mutex.RUnlock()
				}
			} else {
				store.repository.mutex.RLock()
				if err := store.repository.db.Unscoped().First(&Data{}, tt.id).Error; err == nil {
					t.Errorf("still found after delete")
					store.repository.mutex.RUnlock()
					return
				} else {
					store.repository.mutex.RUnlock()
				}
			}

		})
	}
}

func TestGenericStore_DeleteMany(t *testing.T) {
	repository, err := setup(t)
	if err != nil {
		t.Errorf("failed to setup repository: %v", err)
		return
	}
	defer repository.Stop(context.Background())

	store := NewGenericStore[Data](repository)
	createTestData(t, store)

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
			got, err := store.DeleteMany(context.Background(), nil, tt.permanent, tt.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("store.DeleteMany() = %v, want %v", got, tt.want)
			}

			if !tt.permanent {
				store.repository.mutex.RLock()
				if out, err := store.GetMany(context.Background(), nil, tt.ids); err == nil && len(out) != 0 {
					t.Errorf("still found after delete")
					store.repository.mutex.RUnlock()
					return
				} else {
					store.repository.mutex.RUnlock()
				}
			} else {
				store.repository.mutex.RLock()
				var data []*Data
				if err := store.repository.db.Unscoped().Find(&data, tt.ids).Error; err == nil && len(data) != 0 {
					t.Errorf("still found after delete")
					store.repository.mutex.RUnlock()
					return
				} else {
					store.repository.mutex.RUnlock()
				}
			}

		})
	}
}
