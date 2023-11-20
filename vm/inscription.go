package vm

// Content
// only for brc20 current
//
type Content struct {
	Proto     string `json:"p,omitempty"`
	Operation string `json:"op,omitempty"`
	Tick      string `json:"tick,omitempty"`
	Max       string `json:"max,omitempty"`
	Amount    string `json:"amt,omitempty"`
	Limit     string `json:"lim,omitempty"` // option
	Decimal   string `json:"dec,omitempty"` // option
}
