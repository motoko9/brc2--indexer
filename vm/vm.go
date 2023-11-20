package vm

import (
	"encoding/json"
	"github.com/motoko9/model"
	"github.com/motoko9/state"
	"math"
	"strconv"
	"strings"
)

type Vm struct {
}

func New() *Vm {
	vm := &Vm{}
	return vm
}

func (v *Vm) Execute(s *state.State, transaction *model.Transaction) *model.Receipt {
	receipt := &model.Receipt{
		Hash:          transaction.Hash,
		InscriptionId: transaction.InscriptionId,
		Status:        0,
		Msg:           "",
	}
	//
	if transaction.Inscription.ContentType != "application/json" {
		receipt.Status = 0
		receipt.Msg = "not support"
		return receipt
	}
	//
	var r Content
	err := json.Unmarshal(transaction.Inscription.Content, &r)
	if err != nil {
		receipt.Status = 0
		receipt.Msg = "not support"
		return receipt
	}
	//
	if r.Proto != "brc-20" {
		receipt.Status = 0
		receipt.Msg = "not support"
		return receipt
	}
	//
	if r.Tick != "ordi" {
		receipt.Status = 0
		receipt.Msg = "not support"
		return receipt
	}
	s.Reload(strings.ToLower(r.Tick))
	v.ExecuteBrc20(&transaction.Input, &transaction.Output, &r, s, receipt)
	return receipt
}

func (v *Vm) ExecuteBrc20(input *model.Input, output *model.Output, r *Content, s *state.State, receipt *model.Receipt) {
	switch r.Operation {
	case "deploy":
		v.handleBrc20Deploy(input, output, r, s, receipt)
	case "mint":
		v.handleBrc20Mint(input, output, r, s, receipt)
	case "transfer":
		v.handleBrc20Transfer(input, output, r, s, receipt)
	default:
		return
	}
}

func (v *Vm) handleBrc20Deploy(input *model.Input, output *model.Output, r *Content, s *state.State, receipt *model.Receipt) {
	brc20Ticker := strings.ToLower(r.Tick)
	if !s.IsEmpty() {
		return
	}
	//
	s.Create(brc20Ticker)
	//
	var err error
	if r.Max == "" {
		return
	}
	decimal := int64(1)
	if r.Decimal != "" {
		decimal, err = strconv.ParseInt(r.Decimal, 10, 64)
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
	//
	s.Set("info:name", brc20Ticker)
	s.Set("info:decimal", decimal)
	s.Set("info:max", max)
	s.Set("info:limit", limit)
	s.Set("info:total_supply", int64(0))
	//
	receipt.Status = 1
	event := make([]string, 5)
	event[0] = "brc-20"
	event[1] = brc20Ticker
	event[2] = r.Decimal
	event[3] = r.Max
	event[4] = r.Limit
	receipt.Events = append(receipt.Events, model.Event{
		Id:   "deploy",
		Data: event,
	})
}

func (v *Vm) handleBrc20Mint(input *model.Input, output *model.Output, r *Content, s *state.State, receipt *model.Receipt) {
	if s.IsEmpty() {
		return
	}
	//
	amount, err := strconv.ParseInt(r.Amount, 10, 64)
	if err != nil {
		return
	}
	limit := s.Get("info:limit").(int64)
	if amount > limit {
		return
	}
	max := s.Get("info:max").(int64)
	totalSupply := s.Get("info:total_supply").(int64)
	if totalSupply >= max {
		return
	}
	if totalSupply+amount > max {
		amount = max - totalSupply
	}
	//
	s.Set("info:total_supply", totalSupply+amount)
	balance := s.Get("balance:" + input.Address).(int64)
	s.Set("balance:"+input.Address, balance+amount)
	//
	receipt.Status = 1
	event := make([]string, 2)
	event[0] = input.Address
	event[1] = input.Address
	event[2] = r.Amount
	receipt.Events = append(receipt.Events, model.Event{
		Id:   "transfer",
		Data: event,
	})
}

func (v *Vm) handleBrc20Transfer(input *model.Input, output *model.Output, r *Content, s *state.State, receipt *model.Receipt) {
	if s.IsEmpty() {
		return
	}
	//
	fromAddress := input.Address
	toAddress := output.Address
	amount, err := strconv.ParseInt(r.Amount, 10, 64)
	if err != nil {
		return
	}
	fromBalance := s.Get("balance:" + fromAddress).(int64)
	toBalance := s.Get("balance:" + toAddress).(int64)
	s.Set("balance:"+fromAddress, fromBalance-amount)
	s.Set("balance:"+toAddress, toBalance+amount)
	//
	receipt.Status = 1
	event := make([]string, 2)
	event[0] = fromAddress
	event[1] = toAddress
	event[2] = r.Amount
	receipt.Events = append(receipt.Events, model.Event{
		Id:   "transfer",
		Data: event,
	})
}
