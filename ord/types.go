package ord

type Inscription struct {
	Address           string `json:"address"`
	ContentLength     uint64 `json:"content_length"`
	ContentType       string `json:"content_type"`
	GenesisFee        uint64 `json:"genesis_fee"`
	GenesisHeight     uint64 `json:"genesis_height"`
	InscriptionId     string `json:"inscription_id"`
	InscriptionNumber uint64 `json:"inscription_number"`
	OutputValue       uint64 `json:"output_value"`
	SatPoint          string `json:"satpoint"`
	Timestamp         uint64 `json:"timestamp"`
}

type Input struct {
	Id string
}

type Output struct {
	Id      string
	Value   string
	Address string
}

type Transaction struct {
	Hash    string   `json:"hash"`
	Inputs  []Input  `json:"input"`
	Outputs []Output `json:"output"`
}
