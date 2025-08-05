package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToTimepb(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func FromTimepb(t *timestamppb.Timestamp) time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.AsTime()
}
