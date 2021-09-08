package types

import (
	"time"
)

type Transaction struct {
	Type      string     `bson:"type"`
	SubType   string     `bson:"subtype"`
	Time      time.Time  `bson:"time"`
	ID        string     `bson:"id"`
	AccountID string     `bson:"account-id"`
	Fill      *Fill      `bson:"fill"`
	Movements []Movement `bson:"movements"`
}

type Fill struct {
	SecurityID string  `bson:"security-id"`
	Price      float64 `bson:"price"`
	Quantity   float64 `bson:"quantity"`
}

type Movement struct {
	Reason   int32   `bson:"reason"`
	AssetID  uint32  `bson:"asset-id"`
	Quantity float64 `bson:"quantity"`
}
