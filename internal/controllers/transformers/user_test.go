package transformers

import (
	"reflect"
	"testing"
	"time"

	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/libs/utils"
	pb "github.com/tuantran1810/go-di-template/pkg/go_di_template/v1"
)

func TestPbUserTransformer_ToEntity(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()

	tr := &pbUserTransformer{}

	tests := []struct {
		name    string
		user    *pb.User
		want    *entities.User
		wantErr bool
	}{
		{
			name: "success",
			user: &pb.User{
				Id:        uint32(1),
				CreatedAt: utils.ToTimepb(now),
				UpdatedAt: utils.ToTimepb(now),
				Username:  "test",
				Password:  "test",
				Uuid:      "test",
				Name:      "test",
				Email:     utils.Pointer("test@test.com"),
			},
			want: &entities.User{
				ID:        uint(1),
				CreatedAt: now,
				UpdatedAt: now,
				Username:  "test",
				Password:  "test",
				Uuid:      "test",
				Name:      "test",
				Email:     utils.Pointer("test@test.com"),
			},
			wantErr: false,
		},
		{
			name:    "nil input",
			user:    nil,
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tr.ToEntity(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("PbUserTransformer.ToEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PbUserTransformer.ToEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPbUserTransformer_FromEntity(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()

	tr := &pbUserTransformer{}

	tests := []struct {
		name    string
		tr      *pbUserTransformer
		user    *entities.User
		want    *pb.User
		wantErr bool
	}{
		{
			name: "success",
			user: &entities.User{
				ID:        uint(1),
				CreatedAt: now,
				UpdatedAt: now,
				Username:  "test",
				Password:  "test",
				Uuid:      "test",
				Name:      "test",
				Email:     utils.Pointer("test@test.com"),
			},
			want: &pb.User{
				Id:        uint32(1),
				CreatedAt: utils.ToTimepb(now),
				UpdatedAt: utils.ToTimepb(now),
				Username:  "test",
				Password:  "test",
				Uuid:      "test",
				Name:      "test",
				Email:     utils.Pointer("test@test.com"),
			},
			wantErr: false,
		},
		{
			name:    "nil input",
			user:    nil,
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tr.FromEntity(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("PbUserTransformer.FromEntity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PbUserTransformer.FromEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}
