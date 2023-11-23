package model

// Content
// only for brc20 current
//
type Content struct {
	Proto     string `json:"p,omitempty"`
	Operation string `json:"op,omitempty"`
	Name      string `json:"tick,omitempty"`
	Maximum   string `json:"max,omitempty"`
	Amount    string `json:"amt,omitempty"`
	Limit     string `json:"lim,omitempty"` // option
	Decimal   string `json:"dec,omitempty"` // option
}

type Inscription struct {
	Address           string `json:"address"`
	ContentLength     uint64 `json:"content_length"`
	ContentType       string `json:"content_type"`
	Content           []byte `json:"content"`
	GenesisFee        uint64 `json:"genesis_fee"`
	GenesisHeight     uint64 `json:"genesis_height"`
	InscriptionId     string `json:"inscription_id"`
	InscriptionNumber int64  `json:"inscription_number"`
	OutputValue       uint64 `json:"output_value"`
	SatPoint          string `json:"satpoint"`
	Timestamp         uint64 `json:"timestamp"`
}

type Output struct {
	Hash    string
	N       int64
	Address string
	Value   uint64
}

type Transaction struct {
	Output      Output
	Inscription Inscription
}

type Event struct {
	Name string
	Id   string
	Data []string
}

type Context struct {
	Output      Output      `json:"output"`
	Inscription Inscription `json:"inscription"`
	Content     Content
	Status      int    `json:"status"`
	Msg         string `json:"msg"`
	Event       Event  `json:"events"`
}
