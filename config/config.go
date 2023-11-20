package config

type OrdNode struct {
	Url string `json:"url"`
}

type BitcoinNode struct {
	Url  string `json:"url"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type Postgres struct {
	Connect string
}

type Config struct {
	OrdNode     OrdNode     `json:"ord_name"`
	BitcoinNode BitcoinNode `json:"bitcoin_node"`
	Postgres    Postgres    `json:"postgres"`
}
