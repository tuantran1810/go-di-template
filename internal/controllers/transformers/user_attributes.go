package transformers

import (
	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/libs/utils"
	pb "github.com/tuantran1810/go-di-template/pkg/go_di_template/v1"
)

type pbUserAttributesTransformer struct{}
type PbUserAttributesTransformer = entities.ExtendedDataTransformer[pb.UserAttribute, entities.UserAttribute]

func NewPbUserAttributesTransformer() *PbUserAttributesTransformer {
	return entities.NewExtendedDataTransformer(&pbUserAttributesTransformer{})
}

func (t *pbUserAttributesTransformer) ToEntity(data *pb.UserAttribute) (*entities.UserAttribute, error) {
	if data == nil {
		return nil, nil
	}

	return &entities.UserAttribute{
		ID:        uint(data.Id),
		CreatedAt: utils.FromTimepb(data.CreatedAt),
		UpdatedAt: utils.FromTimepb(data.UpdatedAt),
		UserID:    uint(data.UserId),
		Key:       data.Key,
		Value:     data.Value,
	}, nil
}

func (t *pbUserAttributesTransformer) FromEntity(entity *entities.UserAttribute) (*pb.UserAttribute, error) {
	if entity == nil {
		return nil, nil
	}

	return &pb.UserAttribute{
		Id:        uint32(entity.ID),
		CreatedAt: utils.ToTimepb(entity.CreatedAt),
		UpdatedAt: utils.ToTimepb(entity.UpdatedAt),
		UserId:    uint32(entity.UserID),
		Key:       entity.Key,
		Value:     entity.Value,
	}, nil
}
