package processor

import (
	"context"
	"time"

	"github.com/tuantran1810/go-di-template/internal/entities"
)

type Sender interface {
	Send(ctx context.Context, msg *entities.Message) error
}

type DataRepository interface {
	RunTx(ctx context.Context, data any, funcs ...entities.DBTxHandleFunc) (any, error)
}

type DataWriter interface {
	Create(ctx context.Context, tx entities.Transaction, data *entities.Message) (*entities.Message, error)
}

type Processor struct {
	sender Sender
	writer DataWriter
	repo   DataRepository
}

func NewProcessor(sender Sender, writer DataWriter, repo DataRepository) *Processor {
	return &Processor{
		sender: sender,
		writer: writer,
		repo:   repo,
	}
}

func (p *Processor) Start(_ context.Context) error {
	return nil
}

func (p *Processor) Stop(_ context.Context) error {
	return nil
}

func (p *Processor) Process(ctx context.Context, msg *entities.Message) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	_, err := p.repo.RunTx(
		timeoutCtx,
		nil,
		func(ictx context.Context, tx entities.Transaction, data any) (out any, cont bool, err error) {
			if _, ierr := p.writer.Create(ictx, tx, msg); ierr != nil {
				return nil, false, ierr
			}
			return nil, true, nil
		},
		func(ictx context.Context, tx entities.Transaction, data any) (out any, cont bool, err error) {
			return nil, true, p.sender.Send(ictx, msg)
		},
	)
	return err
}
