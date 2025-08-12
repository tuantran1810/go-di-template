package usecases

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/tuantran1810/go-di-template/internal/entities"
	mockUsecases "github.com/tuantran1810/go-di-template/mocks/usecases"
)

func TestLoggingWorker_Inject(t *testing.T) {
	t.Parallel()
	capacity := 100

	w := &LoggingWorker{
		LoggingWorkerConfig: LoggingWorkerConfig{
			BufferCapacity: capacity,
			FlushInterval:  1 * time.Second,
		},
		signalChan: make(chan struct{}, 1),
	}

	t.Run("inject to trigger flush", func(t *testing.T) {
		for i := 0; i < capacity; i++ {
			w.Inject(entities.Message{
				Key:   "test",
				Value: "test",
			})
		}

		<-w.signalChan
	})
}

func TestLoggingWorker_StartInjectAndFlushOnStop(t *testing.T) {
	capacity := 100

	mockMessageRepository := mockUsecases.NewMockIMessageRepository(t)
	mockMessageRepository.EXPECT().
		CreateMany(mock.Anything, mock.Anything, mock.Anything).
		Return([]entities.Message{}, nil)

	w := NewLoggingWorker(
		LoggingWorkerConfig{
			BufferCapacity: capacity,
			FlushInterval:  1 * time.Second,
		},
		mockMessageRepository,
		nil,
	)

	t.Run("start, inject and flush on stop", func(t *testing.T) {
		if err := w.Start(context.Background()); err != nil {
			t.Errorf("failed to start logging worker: %v", err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())

		for i := 0; i < 10; i++ {
			go func() {
				ticker := time.NewTicker(10 * time.Millisecond)
				for {
					select {
					case <-ticker.C:
						w.Inject(entities.Message{
							Key:   "test",
							Value: "test",
						})
					case <-ctx.Done():
						return
					}
				}
			}()
		}

		time.Sleep(1010 * time.Millisecond)
		cancel()

		w.Inject(entities.Message{
			Key:   "test",
			Value: "test",
		})

		if err := w.Stop(context.Background()); err != nil {
			t.Errorf("failed to stop logging worker: %v", err)
		}
	})
}

func TestLoggingWorker_StartInjectAndFlushOnInterval(t *testing.T) {
	capacity := 100000

	mockMessageRepository := mockUsecases.NewMockIMessageRepository(t)
	mockMessageRepository.EXPECT().
		CreateMany(mock.Anything, mock.Anything, mock.Anything).
		Return([]entities.Message{}, nil)

	w := NewLoggingWorker(
		LoggingWorkerConfig{
			BufferCapacity: capacity,
			FlushInterval:  10 * time.Millisecond,
		},
		mockMessageRepository,
		nil,
	)

	t.Run("start, inject and flush on interval", func(t *testing.T) {
		if err := w.Start(context.Background()); err != nil {
			t.Errorf("failed to start logging worker: %v", err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			ticker := time.NewTicker(1 * time.Millisecond)
			for {
				select {
				case <-ticker.C:
					w.Inject(entities.Message{
						Key:   "test",
						Value: "test",
					})
				case <-ctx.Done():
					return
				}
			}
		}()

		time.Sleep(102 * time.Millisecond)
		cancel()

		if err := w.Stop(context.Background()); err != nil {
			t.Errorf("failed to stop logging worker: %v", err)
		}
	})
}

func TestLoggingWorker_flush(t *testing.T) {
	t.Parallel()

	mockMessageRepository := mockUsecases.NewMockIMessageRepository(t)
	mockMessageRepository.EXPECT().
		CreateMany(mock.Anything, mock.Anything, []entities.Message{
			{
				Key:   "test1",
				Value: "test1",
			},
			{
				Key:   "test2",
				Value: "test2",
			},
		}).
		Return([]entities.Message{}, nil)

	mockMessageRepository.EXPECT().
		CreateMany(mock.Anything, mock.Anything, []entities.Message{
			{
				Key:   "test_failed_1",
				Value: "test_failed_1",
			},
			{
				Key:   "test_failed_2",
				Value: "test_failed_2",
			},
		}).
		Return(nil, errors.New("fake error"))

	w := &LoggingWorker{
		lock:              sync.Mutex{},
		messageRepository: mockMessageRepository,
	}

	tests := []struct {
		name     string
		buffer   []entities.Message
		expected []entities.Message
		wantErr  bool
	}{
		{
			name: "no error, buffer is not empty",
			buffer: []entities.Message{
				{
					Key:   "test1",
					Value: "test1",
				},
				{
					Key:   "test2",
					Value: "test2",
				},
			},
			expected: []entities.Message{},
			wantErr:  false,
		},
		{
			name: "error, buffer is not empty",
			buffer: []entities.Message{
				{
					Key:   "test_failed_1",
					Value: "test_failed_1",
				},
				{
					Key:   "test_failed_2",
					Value: "test_failed_2",
				},
			},
			expected: []entities.Message{
				{
					Key:   "test_failed_1",
					Value: "test_failed_1",
				},
				{
					Key:   "test_failed_2",
					Value: "test_failed_2",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w.buffer = tt.buffer
			if err := w.flush(); (err != nil) != tt.wantErr {
				t.Errorf("LoggingWorker.flush() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(w.buffer, tt.expected) {
				t.Errorf("LoggingWorker.flush() = %v, want %v", w.buffer, tt.expected)
			}
		})
	}
}

func TestLoggingWorker_LogAndSend(t *testing.T) {
	t.Parallel()

	mockClient := mockUsecases.NewMockIClient(t)
	mockClient.EXPECT().
		Send(mock.Anything, &entities.Message{
			Key:   "test",
			Value: "test",
		}).
		Return(nil)

	mockClient.EXPECT().
		Send(mock.Anything, &entities.Message{
			Key:   "test_failed",
			Value: "test_failed",
		}).
		Return(errors.New("fake error"))

	w := &LoggingWorker{
		LoggingWorkerConfig: LoggingWorkerConfig{
			BufferCapacity: 100,
			FlushInterval:  1 * time.Second,
		},
		lock:   sync.Mutex{},
		buffer: []entities.Message{},
		client: mockClient,
	}
	tests := []struct {
		name    string
		msg     entities.Message
		wantErr bool
	}{
		{
			name: "no error",
			msg: entities.Message{
				Key:   "test",
				Value: "test",
			},
			wantErr: false,
		},
		{
			name: "error",
			msg: entities.Message{
				Key:   "test_failed",
				Value: "test_failed",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := w.LogAndSend(context.TODO(), tt.msg); (err != nil) != tt.wantErr {
				t.Errorf("LoggingWorker.LogAndSend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
