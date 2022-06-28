package types

import (
	"time"
)

type Account struct {
	ID         uint   `gorm:"primarykey"`
	Name       string `gorm:"unique"`
	ExchangeID uint32
}

type Transaction struct {
	ID          uint `gorm:"primarykey"`
	Type        string
	SubType     string
	Time        time.Time
	ExecutionID string     `gorm:"unique"`
	AccountID   uint       `gorm:"index"`
	Account     Account    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Fill        *Fill      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Movements   []Movement `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Fill struct {
	ID            uint `gorm:"primarykey"`
	Transaction   Transaction
	TransactionID uint    `gorm:"index"`
	Account       Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AccountID     uint    `gorm:"index"`
	SecurityID    int64
	Price         float64
	Quantity      float64
}

type Movement struct {
	ID            uint        `gorm:"primarykey"`
	Transaction   Transaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TransactionID uint        `gorm:"index"`
	Account       Account     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AccountID     uint        `gorm:"index"`
	Reason        int32
	AssetID       uint32
	Quantity      float64
}

type MongoTransaction struct {
	Type      string          `bson:"type"`
	SubType   string          `bson:"subtype"`
	Time      time.Time       `bson:"time"`
	ID        string          `bson:"id"`
	Account   string          `bson:"account"`
	Fill      *MongoFill      `bson:"fill"`
	Movements []MongoMovement `bson:"movements"`
}

type MongoFill struct {
	SecurityID string  `bson:"security-id"`
	Price      float64 `bson:"price"`
	Quantity   float64 `bson:"quantity"`
}

type MongoMovement struct {
	Reason   int32   `bson:"reason"`
	AssetID  uint32  `bson:"asset-id"`
	Quantity float64 `bson:"quantity"`
}
