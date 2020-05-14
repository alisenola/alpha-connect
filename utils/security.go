package utils

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"sort"
)

func SecurityID(typ, symbol, exchange string) uint64 {
	tags := map[string]string{
		"type":     typ,
		"symbol":   symbol,
		"exchange": exchange,
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
