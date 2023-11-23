package db

type Sync struct {
	Id           int64 `gorm:"primaryKey;autoIncrement"`
	SyncHeight   int64
	CommitHeight int64
}

type Transaction struct {
	Hash          string `gorm:"primaryKey"`
	N             int    `gorm:"primaryKey"`
	InscriptionId string `gorm:"primaryKey"`
	Timestamp     int64
	Height        int64
	Receipts      []Receipt `gorm:"foreignKey:Hash,N,InscriptionId;references:Hash,N,InscriptionId"`
}

type Inscription struct {
	InscriptionId     string `gorm:"primaryKey"`
	ContentLength     uint64
	ContentType       string
	Content           []byte
	InscriptionNumber int64
	Owner             string
}

type Receipt struct {
	Hash          string `gorm:"primaryKey"`
	N             int    `gorm:"primaryKey"`
	InscriptionId string `gorm:"primaryKey"`
	Status        int
	Msg           string
	// event
	Data1 string
	Data2 string
	Data3 string
	Data4 string
	Data5 string
	//Transaction Transaction `gorm:"foreignKey:hash,n,inscription_id;references:hash,n,inscription_id"`
}

type Brc20Info struct {
	Name        string `gorm:"primaryKey"`
	Maximum     int64
	Limit       int64
	Decimal     int64
	TotalSupply int64
}

type Brc20Balance struct {
	Name    string `gorm:"primaryKey"`
	Address string `gorm:"primaryKey"`
	Balance int64
}
