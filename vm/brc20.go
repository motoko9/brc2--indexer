package vm

type BRC20Content struct {
	Proto        string `json:"p,omitempty"`
	Operation    string `json:"op,omitempty"`
	Tick         string `json:"tick,omitempty"`
	Max          string `json:"max,omitempty"`
	BRC20Amount  string `json:"amt,omitempty"`
	Limit        string `json:"lim,omitempty"` // option
	BRC20Decimal string `json:"dec,omitempty"` // option
}

type Info struct {
	Max         int64
	Limit       int64
	Decimal     int64
	TotalSupply int64
}

type Brc20 struct {
	Tick     string
	Deployer string
	info     Info
	balances map[string]int64
}
