package stores

import "github.com/tuantran1810/go-di-template/uberfx/internal/models"

type MessageStore struct {
	*GenericStore[models.Message]
}

func NewMessageStore(repository *Repository) *MessageStore {
	return &MessageStore{
		GenericStore: NewGenericStore[models.Message](repository),
	}
}
