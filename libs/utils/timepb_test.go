package utils

import (
	"reflect"
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestToTimepb(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	want := &timestamppb.Timestamp{
		Seconds: now.Unix(),
		Nanos:   int32(now.Nanosecond()),
	}
	t.Run("now", func(t *testing.T) {
		t.Parallel()
		if got := ToTimepb(now); !reflect.DeepEqual(got, want) {
			t.Errorf("ToTimepb() = %v, want %v", got, want)
		}
	})
}

func TestFromTimepb(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	pbNow := ToTimepb(now)
	t.Run("now", func(t *testing.T) {
		t.Parallel()
		if got := FromTimepb(pbNow); !reflect.DeepEqual(got, now) {
			t.Errorf("FromTimepb() = %v, want %v", got, now)
		}
	})

	var blank time.Time
	t.Run("nil input", func(t *testing.T) {
		t.Parallel()
		if got := FromTimepb(nil); !reflect.DeepEqual(got, blank) {
			t.Errorf("FromTimepb() = %v, want %v", got, blank)
		}
	})
}
