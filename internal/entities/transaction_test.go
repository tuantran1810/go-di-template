package entities

import (
	"reflect"
	"testing"

	"gorm.io/gorm"
)

func TestGormTransaction_GetTransaction(t *testing.T) {
	t.Parallel()
	var nilptr *gorm.DB
	tests := []struct {
		tx   *gorm.DB
		want any
	}{
		{
			tx:   &gorm.DB{},
			want: &gorm.DB{},
		},
		{
			tx:   nil,
			want: nilptr,
		},
	}
	for _, tt := range tests {
		t.Run("TestGormTransaction_GetTransaction", func(t *testing.T) {
			t.Parallel()
			tr := &GormTransaction{
				Tx: tt.tx,
			}
			if got := tr.GetTransaction(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GormTransaction.GetTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}
