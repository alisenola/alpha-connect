package utils

import (
	"github.com/gogo/protobuf/types"
)

const FIX_UTCDateOnly_LAYOUT = "20060102"
const FIX_UTCTimeOnly_LAYOUT = "15:04:05.0000"
const FIX_UTCDateTime_LAYOUT = FIX_UTCDateOnly_LAYOUT + "-" + FIX_UTCTimeOnly_LAYOUT

func TimestampToMilli(ts *types.Timestamp) uint64 {
	return uint64(ts.Seconds)*1000 + uint64(ts.Nanos/1000000)
}

func SecondToTimestamp(ts uint64) *types.Timestamp {
	seconds := int64(ts)

	return &types.Timestamp{Seconds: seconds, Nanos: 0}
}

func MilliToTimestamp(ts uint64) *types.Timestamp {
	seconds := int64(ts) / 1000
	nanos := int32((int64(ts) - (seconds * 1000)) * 1000000)

	return &types.Timestamp{Seconds: seconds, Nanos: nanos}
}

func TimestampToMicro(ts *types.Timestamp) uint64 {
	return uint64(ts.Seconds)*1000000 + uint64(ts.Nanos/1000)
}

func MicroToTimestamp(ts uint64) *types.Timestamp {
	seconds := int64(ts) / 1000000
	nanos := int32((int64(ts) - (seconds * 1000000)) * 1000)

	return &types.Timestamp{Seconds: seconds, Nanos: nanos}
}
