package utils

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sort"
)

func SecurityID(typ, symbol, exchange string, maturityDate *timestamppb.Timestamp) uint64 {
	tags := map[string]string{
		"type":     typ,
		"symbol":   symbol,
		"exchange": exchange,
	}
	if maturityDate != nil {
		ts := maturityDate.AsTime()
		tags["maturityDate"] = ts.String()
	}
	var keys []string
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	str := ""
	for _, k := range keys {
		str += fmt.Sprintf("&%s=%s", k, tags[k])
	}
	hashBytes := md5.Sum([]byte(str))
	ID := binary.LittleEndian.Uint64(hashBytes[:8])
	return ID
}
