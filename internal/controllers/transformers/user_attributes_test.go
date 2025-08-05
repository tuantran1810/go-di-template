package transformers

import (
	"reflect"
	"testing"
	"time"

	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/libs/utils"
	pb "github.com/tuantran1810/go-di-template/pkg/go_di_template/v1"
)

func TestPbUserAttributesTransformer_ToEntity(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	tr := &pbUserAttributesTransformer{}

	tests := []struct {
		name    string
		data    *pb.UserAttribute
		want    *entities.UserAttribute
		wantErr bool
	}{
		{
			name: "success",
			data: &pb.UserAttribute{
				Id:        uint32(1),
				CreatedAt: utils.ToTimepb(now),
				UpdatedAt: utils.ToTimepb(now),
				UserId:    uint32(1),
				Key:       "test",
				Value:     "test",
			},
			want: &entities.UserAttribute{
				ID:        uint(1),
				CreatedAt: now,
				UpdatedAt: now,
				UserID:    uint(1),
				Key:       "test",
				Value:     "test",
			},
			wantErr: false,
		},
		{
			name:    "nil input",
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
				t.Errorf("PbUserAttributesTransformer.ToEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PbUserAttributesTransformer.ToEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPbUserAttributesTransformer_FromEntity(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	tr := &pbUserAttributesTransformer{}

	tests := []struct {
		name    string
		entity  *entities.UserAttribute
		want    *pb.UserAttribute
		wantErr bool
	}{
		{
			name: "success",
			entity: &entities.UserAttribute{
				ID:        uint(1),
				CreatedAt: now,
				UpdatedAt: now,
				UserID:    uint(1),
				Key:       "test",
				Value:     "test",
			},
			want: &pb.UserAttribute{
				Id:        uint32(1),
				CreatedAt: utils.ToTimepb(now),
				UpdatedAt: utils.ToTimepb(now),
				UserId:    uint32(1),
				Key:       "test",
				Value:     "test",
			},
			wantErr: false,
		},
		{
			name:    "nil input",
			entity:  nil,
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tr.FromEntity(tt.entity)
			if (err != nil) != tt.wantErr {
				t.Errorf("PbUserAttributesTransformer.FromEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PbUserAttributesTransformer.FromEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}
