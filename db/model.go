package db

type Brc20Transaction struct {
	Hash          string `json:"hash"`
	InscriptionId string `json:"inscription_id"`
	ContentLength uint64 `json:"content_length"`
	ContentType   string `json:"content_type"`
	Content       []byte `json:"content"`
	Timestamp     uint64 `json:"timestamp"`
	Height        uint64 `json:"height"`
	Status        int    `json:"status"`
	Msg           string `json:"msg"`
}

type Brc20Event struct {
	Brc20 string
	Id    string
	Data  []byte
}

type Brc20Info struct {
	Tick    string
	Max     int64
	Limit   int64
	Decimal int64
}

type Brc20Balance struct {
	Tick    string
	Address string
	Balance int64
}

type Inscription struct {
	Name  string
	Proto string
}
