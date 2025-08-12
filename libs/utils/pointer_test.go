package utils

import (
	"reflect"
	"testing"
	"time"
)

func TestPointer(t *testing.T) {
	t.Parallel()

	var vInt int64 = 1
	t.Run("int value", func(t *testing.T) {
		t.Parallel()
		if got := Pointer(vInt); !reflect.DeepEqual(*got, vInt) {
			t.Errorf("Pointer() = %v, want %v", got, vInt)
		}
	})

	var vString = "test"
	t.Run("string value", func(t *testing.T) {
		t.Parallel()
		if got := Pointer(vString); !reflect.DeepEqual(*got, vString) {
			t.Errorf("Pointer() = %v, want %v", got, vString)
		}
	})

	var vBool = true
	t.Run("bool value", func(t *testing.T) {
		t.Parallel()
		if got := Pointer(vBool); !reflect.DeepEqual(*got, vBool) {
			t.Errorf("Pointer() = %v, want %v", got, vBool)
		}
	})

	var vFloat = 1.0
	t.Run("float value", func(t *testing.T) {
		t.Parallel()
		if got := Pointer(vFloat); !reflect.DeepEqual(*got, vFloat) {
			t.Errorf("Pointer() = %v, want %v", got, vFloat)
		}
	})

	var vTime = time.Now()
	t.Run("time value", func(t *testing.T) {
		t.Parallel()
		if got := Pointer(vTime); !reflect.DeepEqual(*got, vTime) {
			t.Errorf("Pointer() = %v, want %v", got, vTime)
		}
	})
}
