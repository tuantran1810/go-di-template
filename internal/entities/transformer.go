package entities

type DataTransformer[T, E any] interface {
	ToEntity(data *T) (entity *E, err error)
	FromEntity(entity *E) (data *T, err error)
}

type ExtendedDataTransformer[T, E any] struct {
	DataTransformer[T, E]
}

// FromEntityArray_I2I: From Entity instance Array to Data instance Array
func (t *ExtendedDataTransformer[T, E]) FromEntityArray_I2I(entities []E) (data []T, err error) {
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

// FromEntityArray_I2P: From Entity instance Array to Data pointer Array
func (t *ExtendedDataTransformer[T, E]) FromEntityArray_I2P(entities []E) (data []*T, err error) {
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

// FromEntityArray_P2I: From Entity pointer Array to Data instance Array
func (t *ExtendedDataTransformer[T, E]) FromEntityArray_P2I(entities []*E) (data []T, err error) {
	if entities == nil {
		return nil, nil
	}

	dataArray := make([]T, 0, len(entities))
	for _, entity := range entities {
		data, err := t.FromEntity(entity)
		if err != nil {
			return nil, err
		}
		dataArray = append(dataArray, *data)
	}

	return dataArray, nil
}

// FromEntityArray_P2P: From Entity pointer Array to Data pointer Array
func (t *ExtendedDataTransformer[T, E]) FromEntityArray_P2P(entities []*E) (data []*T, err error) {
	if entities == nil {
		return nil, nil
	}

	dataArray := make([]*T, 0, len(entities))
	for _, entity := range entities {
		data, err := t.FromEntity(entity)
		if err != nil {
			return nil, err
		}
		dataArray = append(dataArray, data)
	}

	return dataArray, nil
}

// ToEntityArray_I2I: From Data instance Array to Entity instance Array
func (t *ExtendedDataTransformer[T, E]) ToEntityArray_I2I(dataArray []T) (entity []E, err error) {
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

// ToEntityArray_I2P: From Entity instance Array to Data pointer Array
func (t *ExtendedDataTransformer[T, E]) ToEntityArray_I2P(dataArray []T) (entity []*E, err error) {
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

// ToEntityArray_P2I: From Data pointer Array to Entity instance Array
func (t *ExtendedDataTransformer[T, E]) ToEntityArray_P2I(dataArray []*T) (entity []E, err error) {
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

// ToEntityArray_P2P: From Data pointer Array to Entity pointer Array
func (t *ExtendedDataTransformer[T, E]) ToEntityArray_P2P(dataArray []*T) (entity []*E, err error) {
	if dataArray == nil {
		return nil, nil
	}

	entityArray := make([]*E, 0, len(dataArray))
	for _, data := range dataArray {
		entity, err := t.ToEntity(data)
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
