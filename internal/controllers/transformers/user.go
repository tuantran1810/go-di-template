package transformers

import (
	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/libs/utils"
	pb "github.com/tuantran1810/go-di-template/pkg/go_di_template/v1"
)

type pbUserTransformer struct{}
type PbUserTransformer = entities.ExtendedDataTransformer[pb.User, entities.User]

func NewPbUserTransformer() *PbUserTransformer {
	return entities.NewExtendedDataTransformer(&pbUserTransformer{})
}

func (t *pbUserTransformer) ToEntity(user *pb.User) (*entities.User, error) {
	if user == nil {
		return nil, nil
	}

	return &entities.User{
		ID:        uint(user.Id),
		CreatedAt: utils.FromTimepb(user.CreatedAt),
		UpdatedAt: utils.FromTimepb(user.UpdatedAt),
		Username:  user.Username,
		Password:  user.Password,
		Uuid:      user.Uuid,
		Name:      user.Name,
		Email:     user.Email,
	}, nil
}

func (t *pbUserTransformer) FromEntity(user *entities.User) (*pb.User, error) {
	if user == nil {
		return nil, nil
	}

	return &pb.User{
		Id:        uint32(user.ID),
		CreatedAt: utils.ToTimepb(user.CreatedAt),
		UpdatedAt: utils.ToTimepb(user.UpdatedAt),
		Username:  user.Username,
		Password:  user.Password,
		Uuid:      user.Uuid,
		Name:      user.Name,
		Email:     user.Email,
	}, nil
}
