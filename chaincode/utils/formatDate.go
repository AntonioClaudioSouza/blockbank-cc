package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ReturnDate(date *timestamppb.Timestamp) string {
	return date.AsTime().In(time.FixedZone("UTC-3", -3*60*60)).Format(time.RFC3339)
}
