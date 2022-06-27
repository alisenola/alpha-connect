package utils

import (
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const FIX_UTCDateOnly_LAYOUT = "20060102"
const FIX_UTCTimeOnly_LAYOUT = "15:04:05.0000"
const FIX_UTCDateTime_LAYOUT = FIX_UTCDateOnly_LAYOUT + "-" + FIX_UTCTimeOnly_LAYOUT

func TimestampToMilli(ts *timestamppb.Timestamp) uint64 {
	return uint64(ts.Seconds)*1000 + uint64(ts.Nanos/1000000)
}

func TimestampToSeconds(ts *timestamppb.Timestamp) uint64 {
	return uint64(ts.Seconds)
}

func SecondToTimestamp(ts uint64) *timestamppb.Timestamp {
	seconds := int64(ts)

	return &timestamppb.Timestamp{Seconds: seconds, Nanos: 0}
}

func MilliToTimestamp(ts uint64) *timestamppb.Timestamp {
	seconds := int64(ts) / 1000
	nanos := int32((int64(ts) - (seconds * 1000)) * 1000000)

	return &timestamppb.Timestamp{Seconds: seconds, Nanos: nanos}
}

func TimestampToMicro(ts *timestamppb.Timestamp) uint64 {
	return uint64(ts.Seconds)*1000000 + uint64(ts.Nanos/1000)
}

func MicroToTimestamp(ts uint64) *timestamppb.Timestamp {
	seconds := int64(ts) / 1000000
	nanos := int32((int64(ts) - (seconds * 1000000)) * 1000)

	return &timestamppb.Timestamp{Seconds: seconds, Nanos: nanos}
}

func NanoToTimestamp(ts uint64) *timestamppb.Timestamp {
	seconds := int64(ts) / 1000000000
	nanos := int32(int64(ts) - (seconds * 1000000000))

	return &timestamppb.Timestamp{Seconds: seconds, Nanos: nanos}
}
