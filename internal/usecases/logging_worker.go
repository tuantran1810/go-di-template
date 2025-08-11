package usecases

import (
	"context"
	"sync"
	"time"

	"github.com/tuantran1810/go-di-template/internal/entities"
)

type LoggingWorkerConfig struct {
	BufferCapacity int
	FlushInterval  time.Duration
}

type LoggingWorker struct {
	LoggingWorkerConfig
	lock         sync.Mutex
	messageStore IMessageStore
	buffer       []entities.Message
	cancelCtx    context.Context
	cancelFunc   context.CancelFunc
	signalChan   chan struct{}
}

func NewLoggingWorker(
	config LoggingWorkerConfig,
	store IMessageStore,
) *LoggingWorker {
	cancelCtx, cancel := context.WithCancel(context.Background())
	return &LoggingWorker{
		LoggingWorkerConfig: config,
		lock:                sync.Mutex{},
		messageStore:        store,
		buffer:              make([]entities.Message, 0, config.BufferCapacity),
		cancelCtx:           cancelCtx,
		cancelFunc:          cancel,
		signalChan:          make(chan struct{}, 1),
	}
}

func (w *LoggingWorker) flush() error {
	if len(w.buffer) == 0 {
		return nil
	}

	w.lock.Lock()
	tmp := w.buffer
	w.buffer = make([]entities.Message, 0, w.BufferCapacity)
	w.lock.Unlock()

	_, err := w.messageStore.CreateMany(w.cancelCtx, nil, tmp)
	if err != nil {
		w.lock.Lock()
		w.buffer = append(tmp, w.buffer...)
		w.lock.Unlock()
		return err
	}

	return nil
}

func (w *LoggingWorker) worker() {
	ticker := time.NewTicker(w.FlushInterval)

	for {
		select {
		case <-w.cancelCtx.Done():
			return
		case <-ticker.C:
			if err := w.flush(); err != nil {
				log.Error("failed to flush buffer on interval", "error", err)
			}
		case <-w.signalChan:
			if err := w.flush(); err != nil {
				log.Error("failed to flush buffer on signal", "error", err)
			}
		}
	}
}

func (w *LoggingWorker) Start(ctx context.Context) error {
	go w.worker()
	return nil
}

func (w *LoggingWorker) Stop(ctx context.Context) error {
	close(w.signalChan)
	w.cancelFunc()
	return w.flush()
}

func (w *LoggingWorker) Inject(msg entities.Message) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.buffer = append(w.buffer, msg)
	if len(w.buffer) >= w.BufferCapacity {
		w.signalChan <- struct{}{}
	}
}
