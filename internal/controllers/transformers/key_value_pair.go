package transformers

import (
	"github.com/tuantran1810/go-di-template/internal/entities"
	pb "github.com/tuantran1810/go-di-template/pkg/go_di_template/v1"
)

type pbKeyValuePairTransformer struct{}
type PbKeyValuePairTransformer = entities.ExtendedDataTransformer[pb.KeyValuePair, entities.KeyValuePair]

func NewPbKeyValuePairTransformer() *PbKeyValuePairTransformer {
	return entities.NewExtendedDataTransformer(&pbKeyValuePairTransformer{})
}

func (t *pbKeyValuePairTransformer) ToEntity(data *pb.KeyValuePair) (*entities.KeyValuePair, error) {
	if data == nil {
		return nil, nil
	}

	return &entities.KeyValuePair{
		Key:   data.Key,
		Value: data.Value,
	}, nil
}

func (t *pbKeyValuePairTransformer) FromEntity(entity *entities.KeyValuePair) (*pb.KeyValuePair, error) {
	if entity == nil {
		return nil, nil
	}

	return &pb.KeyValuePair{
		Key:   entity.Key,
		Value: entity.Value,
	}, nil
}
