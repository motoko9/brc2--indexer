package ord

type Inscription struct {
	Address           string `json:"address"`
	ContentLength     uint64 `json:"content_length"`
	ContentType       string `json:"content_type"`
	GenesisFee        uint64 `json:"genesis_fee"`
	GenesisHeight     uint64 `json:"genesis_height"`
	InscriptionId     string `json:"inscription_id"`
	InscriptionNumber int64  `json:"inscription_number"`
	OutputValue       uint64 `json:"output_value"`
	SatPoint          string `json:"satpoint"`
	Timestamp         uint64 `json:"timestamp"`
}

type Output struct {
	Value        uint64   `json:"value"`
	Address      string   `json:"address"`
	Inscriptions []string `json:"inscriptions"`
}
