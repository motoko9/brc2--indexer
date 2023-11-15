package vm

import (
	"encoding/json"
	"github.com/motoko9/model"
	"math"
	"strconv"
	"strings"
)

type Vm struct {
	states map[string]*Brc20
}

func New() *Vm {
	vm := &Vm{
		states: make(map[string]*Brc20),
	}
	return vm
}

func (v *Vm) Execute(transaction *model.Brc20Transaction) {
	//
	if transaction.Inscription.ContentType != "application/json" {
		return
	}
	//
	var r BRC20Content
	err := json.Unmarshal(transaction.Inscription.Content, &r)
	if err != nil {
		return
	}
	//
	if r.Proto != "brc-20" {
		return
	}
	//
	if r.Tick != "ordi" {
		return
	}
	v.ExecuteBrc20(&transaction.Input, &transaction.Output, &r)
}

func (v *Vm) ExecuteBrc20(input *model.Input, output *model.Output, r *BRC20Content) {
	switch r.Operation {
	case "deploy":
		v.handleBrc20Deploy(input, output, r)
	case "mint":
		v.handleBrc20Mint(input, output, r)
	case "transfer":
		v.handleBrc20Transfer(input, output, r)
	default:
		return
	}
}

func (v *Vm) handleBrc20Deploy(input *model.Input, output *model.Output, r *BRC20Content) {
	brc20Ticker := strings.ToLower(r.Tick)
	state, ok := v.states[brc20Ticker]
	if ok {
		return
	}
	//
	var err error
	if r.Max == "" {
		return
	}
	decimal := int64(1)
	if r.BRC20Decimal != "" {
		decimal, err = strconv.ParseInt(r.BRC20Decimal, 10, 64)
		if err != nil {
			return
		}
	}
	max, err := strconv.ParseInt(r.Max, 10, 64)
	if err != nil {
		return
	}
	limit := int64(math.MaxInt64)
	if r.Limit != "" {
		limit, err = strconv.ParseInt(r.Limit, 10, 64)
		if err != nil {
			return
		}
	}
	state = &Brc20{
		Tick:     brc20Ticker,
		Deployer: input.Address,
		info: Info{
			Decimal:     decimal,
			Max:         max,
			Limit:       limit,
			TotalSupply: 0,
		},
		balances: make(map[string]int64),
	}
	v.states[brc20Ticker] = state
}

func (v *Vm) handleBrc20Mint(input *model.Input, output *model.Output, r *BRC20Content) {
	brc20Ticker := strings.ToLower(r.Tick)
	state, ok := v.states[brc20Ticker]
	if !ok {
		return
	}
	//
	amount, err := strconv.ParseInt(r.BRC20Amount, 10, 64)
	if err != nil {
		return
	}
	if amount > state.info.Limit {
		return
	}
	if state.info.TotalSupply >= state.info.Max {
		return
	}
	if state.info.TotalSupply+amount > state.info.Max {
		amount = state.info.Max - state.info.TotalSupply
	}
	//
	state.info.TotalSupply += amount
	state.balances[input.Address] = state.balances[input.Address] + amount
}

func (v *Vm) handleBrc20Transfer(input *model.Input, output *model.Output, r *BRC20Content) {
	brc20Ticker := strings.ToLower(r.Tick)
	state, ok := v.states[brc20Ticker]
	if !ok {
		return
	}
	//
	fromAddress := input.Address
	toAddress := output.Address
	amount, err := strconv.ParseInt(r.BRC20Amount, 10, 64)
	if err != nil {
		return
	}
	if state.balances[fromAddress] < amount {
		return
	}
	state.balances[fromAddress] = state.balances[fromAddress] - amount
	state.balances[toAddress] = state.balances[toAddress] + amount
}
