package utils

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestGetCorrelationID(t *testing.T) {
	t.Parallel()

	// Test with valid correlation ID in context
	expectedID := "test-correlation-id-123"
	ctxWithID := context.WithValue(context.Background(), XCorrelationID, expectedID)

	// Test with invalid type in context (not string)
	ctxWithInvalidType := context.WithValue(context.Background(), XCorrelationID, 12345)

	tests := []struct {
		name    string
		ctx     context.Context
		want    string
		wantGen bool // true if we expect a generated UUID
	}{
		{
			name:    "returns correlation ID from context when present",
			ctx:     ctxWithID,
			want:    expectedID,
			wantGen: false,
		},
		{
			name:    "generates new correlation ID when context is empty",
			ctx:     context.Background(),
			want:    "",
			wantGen: true,
		},
		{
			name:    "generates new correlation ID when context value is nil",
			ctx:     context.WithValue(context.Background(), XCorrelationID, nil),
			want:    "",
			wantGen: true,
		},
		{
			name:    "generates new correlation ID when context value is not string",
			ctx:     ctxWithInvalidType,
			want:    "",
			wantGen: true,
		},
		{
			name:    "generates new correlation ID when context value is empty string",
			ctx:     context.WithValue(context.Background(), XCorrelationID, ""),
			want:    "",
			wantGen: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := GetCorrelationID(tt.ctx)

			if tt.wantGen {
				// Verify it's a valid UUID format
				if _, err := uuid.Parse(got); err != nil {
					t.Errorf("GetCorrelationID() generated invalid UUID: %v, error: %v", got, err)
				}
				if got == "" {
					t.Errorf("GetCorrelationID() generated empty string, expected valid UUID")
				}
			} else {
				if got != tt.want {
					t.Errorf("GetCorrelationID() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGenerateCorrelationID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "generates valid UUID",
		},
		{
			name: "generates unique UUIDs on multiple calls",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.name == "generates valid UUID" {
				got := GenerateCorrelationID()

				// Verify it's not empty
				if got == "" {
					t.Errorf("GenerateCorrelationID() returned empty string")
				}

				// Verify it's a valid UUID format
				if _, err := uuid.Parse(got); err != nil {
					t.Errorf("GenerateCorrelationID() generated invalid UUID: %v, error: %v", got, err)
				}
			}

			if tt.name == "generates unique UUIDs on multiple calls" {
				// Generate multiple IDs and verify they're unique
				ids := make(map[string]bool)
				iterations := 100

				for i := 0; i < iterations; i++ {
					id := GenerateCorrelationID()
					if ids[id] {
						t.Errorf("GenerateCorrelationID() generated duplicate UUID: %v", id)
					}
					ids[id] = true
				}

				if len(ids) != iterations {
					t.Errorf("GenerateCorrelationID() expected %d unique IDs, got %d", iterations, len(ids))
				}
			}
		})
	}
}

func TestInjectCorrelationIDToContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		ctx           context.Context
		correlationID string
		want          string
	}{
		{
			name:          "injects correlation ID into empty context",
			ctx:           context.Background(),
			correlationID: "test-id-123",
			want:          "test-id-123",
		},
		{
			name:          "injects correlation ID into context with existing value",
			ctx:           context.WithValue(context.Background(), XCorrelationID, "old-id"),
			correlationID: "new-id-456",
			want:          "new-id-456",
		},
		{
			name:          "injects empty correlation ID",
			ctx:           context.Background(),
			correlationID: "",
			want:          "",
		},
		{
			name:          "injects UUID correlation ID",
			ctx:           context.Background(),
			correlationID: "550e8400-e29b-41d4-a716-446655440000",
			want:          "550e8400-e29b-41d4-a716-446655440000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			newCtx := InjectCorrelationIDToContext(tt.ctx, tt.correlationID)

			// Verify the new context is not the same as the original
			if newCtx == tt.ctx && tt.correlationID != "" {
				t.Errorf("InjectCorrelationIDToContext() returned same context instance")
			}

			// Verify the correlation ID was injected correctly
			got := newCtx.Value(XCorrelationID)
			if got != tt.want {
				t.Errorf("InjectCorrelationIDToContext() injected value = %v, want %v", got, tt.want)
			}

			// Verify using GetCorrelationID function
			gotFromGet := GetCorrelationID(newCtx)
			if gotFromGet != tt.want {
				t.Errorf("GetCorrelationID() from injected context = %v, want %v", gotFromGet, tt.want)
			}
		})
	}
}
