package entities_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/tuantran1810/go-di-template/internal/entities"
)

type Data struct {
	Key_   string
	Value_ string
}

type DataEntity struct {
	Key   string
	Value string
}

type DataTransformer struct{}

func (t *DataTransformer) ToEntity(data *Data) (*DataEntity, error) {
	if data == nil {
		return nil, nil
	}

	if data.Key_ == "" {
		return nil, fmt.Errorf("%w - input data is empty", entities.ErrInvalid)
	}

	return &DataEntity{
		Key:   data.Key_,
		Value: data.Value_,
	}, nil
}

func (t *DataTransformer) FromEntity(entity *DataEntity) (*Data, error) {
	if entity == nil {
		return nil, nil
	}

	if entity.Key == "" {
		return nil, fmt.Errorf("%w - input entity is empty", entities.ErrInvalid)
	}

	return &Data{
		Key_:   entity.Key,
		Value_: entity.Value,
	}, nil
}

func Test_FromEntityArray_I2I(t *testing.T) {
	t.Parallel()
	eTransformer := entities.NewExtendedDataTransformer(&DataTransformer{})
	tests := []struct {
		name    string
		in      []DataEntity
		want    []Data
		wantErr bool
	}{
		{
			name:    "nil input",
			in:      nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty input",
			in:      []DataEntity{},
			want:    []Data{},
			wantErr: false,
		},
		{
			name: "error, empty entity",
			in: []DataEntity{
				{
					Key:   "",
					Value: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy case",
			in: []DataEntity{
				{
					Key:   "key1",
					Value: "value1",
				},
				{
					Key:   "key2",
					Value: "value2",
				},
			},
			want: []Data{
				{
					Key_:   "key1",
					Value_: "value1",
				},
				{
					Key_:   "key2",
					Value_: "value2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := eTransformer.FromEntityArray_I2I(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromEntityArray_I2I() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromEntityArray_I2I() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FromEntityArray_I2P(t *testing.T) {
	t.Parallel()
	eTransformer := entities.NewExtendedDataTransformer(&DataTransformer{})
	tests := []struct {
		name    string
		in      []DataEntity
		want    []*Data
		wantErr bool
	}{
		{
			name:    "nil input",
			in:      nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty input",
			in:      []DataEntity{},
			want:    []*Data{},
			wantErr: false,
		},
		{
			name: "error, empty entity",
			in: []DataEntity{
				{
					Key:   "",
					Value: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy case",
			in: []DataEntity{
				{
					Key:   "key1",
					Value: "value1",
				},
				{
					Key:   "key2",
					Value: "value2",
				},
			},
			want: []*Data{
				{
					Key_:   "key1",
					Value_: "value1",
				},
				{
					Key_:   "key2",
					Value_: "value2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := eTransformer.FromEntityArray_I2P(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromEntityArray_I2P() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromEntityArray_I2P() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FromEntityArray_P2I(t *testing.T) {
	t.Parallel()
	eTransformer := entities.NewExtendedDataTransformer(&DataTransformer{})
	tests := []struct {
		name    string
		in      []*DataEntity
		want    []Data
		wantErr bool
	}{
		{
			name:    "nil input",
			in:      nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty input",
			in:      []*DataEntity{},
			want:    []Data{},
			wantErr: false,
		},
		{
			name: "error, empty entity",
			in: []*DataEntity{
				{
					Key:   "",
					Value: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy case",
			in: []*DataEntity{
				{
					Key:   "key1",
					Value: "value1",
				},
				{
					Key:   "key2",
					Value: "value2",
				},
			},
			want: []Data{
				{
					Key_:   "key1",
					Value_: "value1",
				},
				{
					Key_:   "key2",
					Value_: "value2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := eTransformer.FromEntityArray_P2I(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromEntityArray_P2I() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromEntityArray_P2I() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FromEntityArray_P2P(t *testing.T) {
	t.Parallel()
	eTransformer := entities.NewExtendedDataTransformer(&DataTransformer{})
	tests := []struct {
		name    string
		in      []*DataEntity
		want    []*Data
		wantErr bool
	}{
		{
			name:    "nil input",
			in:      nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty input",
			in:      []*DataEntity{},
			want:    []*Data{},
			wantErr: false,
		},
		{
			name: "error, empty entity",
			in: []*DataEntity{
				{
					Key:   "",
					Value: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy case",
			in: []*DataEntity{
				{
					Key:   "key1",
					Value: "value1",
				},
				{
					Key:   "key2",
					Value: "value2",
				},
			},
			want: []*Data{
				{
					Key_:   "key1",
					Value_: "value1",
				},
				{
					Key_:   "key2",
					Value_: "value2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := eTransformer.FromEntityArray_P2P(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromEntityArray_P2P() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromEntityArray_P2P() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ToEntityArray_I2I(t *testing.T) {
	t.Parallel()
	eTransformer := entities.NewExtendedDataTransformer(&DataTransformer{})
	tests := []struct {
		name    string
		in      []Data
		want    []DataEntity
		wantErr bool
	}{
		{
			name:    "nil input",
			in:      nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty input",
			in:      []Data{},
			want:    []DataEntity{},
			wantErr: false,
		},
		{
			name: "error, empty entity",
			in: []Data{
				{
					Key_:   "",
					Value_: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy case",
			in: []Data{
				{
					Key_:   "key1",
					Value_: "value1",
				},
				{
					Key_:   "key2",
					Value_: "value2",
				},
			},
			want: []DataEntity{
				{
					Key:   "key1",
					Value: "value1",
				},
				{
					Key:   "key2",
					Value: "value2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := eTransformer.ToEntityArray_I2I(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToEntityArray_I2I() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToEntityArray_I2I() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ToEntityArray_P2I(t *testing.T) {
	t.Parallel()
	eTransformer := entities.NewExtendedDataTransformer(&DataTransformer{})
	tests := []struct {
		name    string
		in      []*Data
		want    []DataEntity
		wantErr bool
	}{
		{
			name:    "nil input",
			in:      nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty input",
			in:      []*Data{},
			want:    []DataEntity{},
			wantErr: false,
		},
		{
			name: "error, empty entity",
			in: []*Data{
				{
					Key_:   "",
					Value_: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy case",
			in: []*Data{
				{
					Key_:   "key1",
					Value_: "value1",
				},
				{
					Key_:   "key2",
					Value_: "value2",
				},
			},
			want: []DataEntity{
				{
					Key:   "key1",
					Value: "value1",
				},
				{
					Key:   "key2",
					Value: "value2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := eTransformer.ToEntityArray_P2I(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToEntityArray_P2I() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToEntityArray_P2I() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ToEntityArray_I2P(t *testing.T) {
	t.Parallel()
	eTransformer := entities.NewExtendedDataTransformer(&DataTransformer{})
	tests := []struct {
		name    string
		in      []Data
		want    []*DataEntity
		wantErr bool
	}{
		{
			name:    "nil input",
			in:      nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty input",
			in:      []Data{},
			want:    []*DataEntity{},
			wantErr: false,
		},
		{
			name: "error, empty entity",
			in: []Data{
				{
					Key_:   "",
					Value_: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy case",
			in: []Data{
				{
					Key_:   "key1",
					Value_: "value1",
				},
				{
					Key_:   "key2",
					Value_: "value2",
				},
			},
			want: []*DataEntity{
				{
					Key:   "key1",
					Value: "value1",
				},
				{
					Key:   "key2",
					Value: "value2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := eTransformer.ToEntityArray_I2P(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToEntityArray_I2P() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToEntityArray_I2P() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ToEntityArray_P2P(t *testing.T) {
	t.Parallel()
	eTransformer := entities.NewExtendedDataTransformer(&DataTransformer{})
	tests := []struct {
		name    string
		in      []*Data
		want    []*DataEntity
		wantErr bool
	}{
		{
			name:    "nil input",
			in:      nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty input",
			in:      []*Data{},
			want:    []*DataEntity{},
			wantErr: false,
		},
		{
			name: "error, empty entity",
			in: []*Data{
				{
					Key_:   "",
					Value_: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "happy case",
			in: []*Data{
				{
					Key_:   "key1",
					Value_: "value1",
				},
				{
					Key_:   "key2",
					Value_: "value2",
				},
			},
			want: []*DataEntity{
				{
					Key:   "key1",
					Value: "value1",
				},
				{
					Key:   "key2",
					Value: "value2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := eTransformer.ToEntityArray_P2P(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToEntityArray_P2P() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToEntityArray_P2P() = %v, want %v", got, tt.want)
			}
		})
	}
}
