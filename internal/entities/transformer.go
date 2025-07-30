package entities

type DataTransformer[T, E any] interface {
	ToEntity(data *T) (entity *E, err error)
	FromEntity(entity *E) (data *T, err error)
}

type ExtendedDataTransformer[T, E any] struct {
	DataTransformer[T, E]
}

func (t *ExtendedDataTransformer[T, E]) FromEntityArray(entity []E) (data []T, err error) {
	if entity == nil {
		return nil, nil
	}

	dataArray := make([]T, 0, len(entity))
	for _, entity := range entity {
		data, err := t.FromEntity(&entity)
		if err != nil {
			return nil, err
		}
		dataArray = append(dataArray, *data)
	}

	return dataArray, nil
}

func (t *ExtendedDataTransformer[T, E]) ToEntityArray(data []T) (entity []E, err error) {
	if data == nil {
		return nil, nil
	}

	entityArray := make([]E, 0, len(data))
	for _, data := range data {
		entity, err := t.ToEntity(&data)
		if err != nil {
			return nil, err
		}
		entityArray = append(entityArray, *entity)
	}

	return entityArray, nil
}

func NewExtendedDataTransformer[T, E any](transformer DataTransformer[T, E]) *ExtendedDataTransformer[T, E] {
	return &ExtendedDataTransformer[T, E]{transformer}
}
