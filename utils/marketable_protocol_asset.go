package utils

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
)

func MarketableProtocolAssetID(protocolAssetID uint64, marketID uint32) uint64 {
	str := fmt.Sprintf("&protocolAssetID=%d", protocolAssetID)
	str += fmt.Sprintf("&marketID=%d", marketID)
	sum := md5.Sum([]byte(str))
	return binary.LittleEndian.Uint64(sum[8:])
}
