package db

import "gorm.io/gorm"

type Brc20Transaction struct {
	gorm.Model    `json:"-"`
	Hash          string `json:"hash"`
	InscriptionId string `json:"inscription_id"`
	ContentLength uint64 `json:"content_length"`
	ContentType   string `json:"content_type"`
	Content       []byte `json:"content"`
	Timestamp     uint64 `json:"timestamp"`
	Height        uint64 `json:"height"`
}

type Brc20Receipt struct {
	gorm.Model    `json:"-"`
	Hash          string `json:"hash"`
	InscriptionId string `json:"inscription_id"`
	Status        int    `json:"status"`
	Msg           string `json:"msg"`
	//Events        []Brc20Event `gorm:"foreignKey:ID" json:"events"`
}

type Brc20Event struct {
	gorm.Model    `json:"-"`
	Brc20         string
	InscriptionId string
	Data1         string
	Data2         string
	Data3         string
	Data4         string
	Data5         string
}

type Brc20Info struct {
	gorm.Model  `json:"-"`
	Name        string
	Maximum     int64
	Limit       int64
	Decimal     int64
	TotalSupply int64
}

type Brc20Balance struct {
	gorm.Model `json:"-"`
	Name       string
	Address    string
	Balance    int64
}

type Inscription struct {
	gorm.Model `json:"-"`
	Name       string
	Proto      string
}
