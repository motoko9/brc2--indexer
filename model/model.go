package model

type Inscription struct {
	Address           string `json:"address"`
	ContentLength     uint64 `json:"content_length"`
	ContentType       string `json:"content_type"`
	Content           []byte `json:"content"`
	GenesisFee        uint64 `json:"genesis_fee"`
	GenesisHeight     uint64 `json:"genesis_height"`
	InscriptionId     string `json:"inscription_id"`
	InscriptionNumber uint64 `json:"inscription_number"`
	OutputValue       uint64 `json:"output_value"`
	SatPoint          string `json:"satpoint"`
	Timestamp         uint64 `json:"timestamp"`
}

type Input struct {
	Hash    string
	N       int64
	Address string
	Value   int64
}

type Output struct {
	N       int64
	Address string
	Value   int64
}

type Transaction struct {
	Hash          string `json:"hash"`
	InscriptionId string `json:"inscription_id"`
	Input         Input  `json:"from"`
	Output        Output `json:"to"`
	Inscription   Inscription
}

type Event struct {
	Id   string
	Data []string
}

type Receipt struct {
	Hash          string  `json:"hash"`
	InscriptionId string  `json:"inscription_id"`
	Status        int     `json:"status"`
	Msg           string  `json:"msg"`
	Events        []Event `json:"events"`
}
