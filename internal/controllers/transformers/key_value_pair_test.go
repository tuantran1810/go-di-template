package transformers

import (
	"reflect"
	"testing"

	"github.com/tuantran1810/go-di-template/internal/entities"
	pb "github.com/tuantran1810/go-di-template/pkg/go_di_template/v1"
)

func Test_pbKeyValuePairTransformer_ToEntity(t *testing.T) {
	t.Parallel()
	tr := NewPbKeyValuePairTransformer()

	tests := []struct {
		name    string
		data    *pb.KeyValuePair
		want    *entities.KeyValuePair
		wantErr bool
	}{
		{
			name: "success",
			data: &pb.KeyValuePair{
				Key:   "key",
				Value: "value",
			},
			want: &entities.KeyValuePair{
				Key:   "key",
				Value: "value",
			},
			wantErr: false,
		},
		{
			name:    "nil",
			data:    nil,
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tr.ToEntity(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("pbKeyValuePairTransformer.ToEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pbKeyValuePairTransformer.ToEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pbKeyValuePairTransformer_FromEntity(t *testing.T) {
	t.Parallel()
	tr := NewPbKeyValuePairTransformer()

	tests := []struct {
		name    string
		entity  *entities.KeyValuePair
		want    *pb.KeyValuePair
		wantErr bool
	}{
		{
			name: "success",
			entity: &entities.KeyValuePair{
				Key:   "key",
				Value: "value",
			},
			want: &pb.KeyValuePair{
				Key:   "key",
				Value: "value",
			},
			wantErr: false,
		},
		{
			name:    "nil",
			entity:  nil,
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tr.FromEntity(tt.entity)
			if (err != nil) != tt.wantErr {
				t.Errorf("pbKeyValuePairTransformer.FromEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pbKeyValuePairTransformer.FromEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}
