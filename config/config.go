package config

type OrdNode struct {
	Url string `json:"url"`
}

type Config struct {
	OrdNode OrdNode `json:"ord_name"`
}
