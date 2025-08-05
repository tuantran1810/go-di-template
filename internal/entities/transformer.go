package entities

type DataTransformer[T, E any] interface {
	ToEntity(data *T) (entity *E, err error)
	FromEntity(entity *E) (data *T, err error)
}

type ExtendedDataTransformer[T, E any] struct {
	DataTransformer[T, E]
}

func (t *ExtendedDataTransformer[T, E]) FromEntityArray(entities []E) (data []T, err error) {
	if entities == nil {
		return nil, nil
	}

	dataArray := make([]T, 0, len(entities))
	for _, entity := range entities {
		data, err := t.FromEntity(&entity)
		if err != nil {
			return nil, err
		}
		dataArray = append(dataArray, *data)
	}

	return dataArray, nil
}

func (t *ExtendedDataTransformer[T, E]) FromEntityToPtrArray(entities []E) (data []*T, err error) {
	if entities == nil {
		return nil, nil
	}

	dataArray := make([]*T, 0, len(entities))
	for _, entity := range entities {
		data, err := t.FromEntity(&entity)
		if err != nil {
			return nil, err
		}
		dataArray = append(dataArray, data)
	}

	return dataArray, nil
}

func (t *ExtendedDataTransformer[T, E]) ToEntityArray(dataArray []T) (entity []E, err error) {
	if dataArray == nil {
		return nil, nil
	}

	entityArray := make([]E, 0, len(dataArray))
	for _, data := range dataArray {
		entity, err := t.ToEntity(&data)
		if err != nil {
			return nil, err
		}
		entityArray = append(entityArray, *entity)
	}

	return entityArray, nil
}

func (t *ExtendedDataTransformer[T, E]) PtrToEntityArray(dataArray []*T) (entity []E, err error) {
	if dataArray == nil {
		return nil, nil
	}

	entityArray := make([]E, 0, len(dataArray))
	for _, data := range dataArray {
		entity, err := t.ToEntity(data)
		if err != nil {
			return nil, err
		}
		entityArray = append(entityArray, *entity)
	}

	return entityArray, nil
}

func (t *ExtendedDataTransformer[T, E]) ToEntityPtrArray(dataArray []T) (entity []*E, err error) {
	if dataArray == nil {
		return nil, nil
	}

	entityArray := make([]*E, 0, len(dataArray))
	for _, data := range dataArray {
		entity, err := t.ToEntity(&data)
		if err != nil {
			return nil, err
		}
		entityArray = append(entityArray, entity)
	}

	return entityArray, nil
}

func NewExtendedDataTransformer[T, E any](transformer DataTransformer[T, E]) *ExtendedDataTransformer[T, E] {
	return &ExtendedDataTransformer[T, E]{transformer}
}
