package vm

import (
	"encoding/json"
	"github.com/hashicorp/go-hclog"
	"github.com/motoko9/model"
	"github.com/motoko9/state"
	"math"
	"strconv"
	"strings"
)

type Vm struct {
	log hclog.Logger
}

func New(log hclog.Logger) *Vm {
	vm := &Vm{
		log: log.Named("vm"),
	}
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
	var r model.Content
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
	if r.Name != "ordi" {
		receipt.Status = 0
		receipt.Msg = "not support"
		return receipt
	}
	s.Load(strings.ToLower(r.Name))
	v.ExecuteBrc20(&transaction.Input, &transaction.Output, &r, s, receipt)
	return receipt
}

func (v *Vm) ExecuteBrc20(input *model.Input, output *model.Output, r *model.Content, s *state.State, receipt *model.Receipt) {
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

func (v *Vm) handleBrc20Deploy(input *model.Input, output *model.Output, r *model.Content, s *state.State, receipt *model.Receipt) {
	name := strings.ToLower(r.Name)
	if !s.IsEmpty(name) {
		return
	}
	//
	s.Create(name)
	//
	var err error
	if r.Maximum == "" {
		return
	}
	decimal := int64(1)
	if r.Decimal != "" {
		decimal, err = strconv.ParseInt(r.Decimal, 10, 64)
		if err != nil {
			return
		}
	}
	maximum, err := strconv.ParseInt(r.Maximum, 10, 64)
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
	s.Set("info:name", name)
	s.Set("info:decimal", decimal)
	s.Set("info:maximum", maximum)
	s.Set("info:limit", limit)
	s.Set("info:total_supply", int64(0))
	//
	receipt.Status = 1
	event := make([]string, 3)
	event[0] = r.Decimal
	event[1] = r.Maximum
	event[2] = r.Limit
	receipt.Events = append(receipt.Events, model.Event{
		Name: name,
		Id:   "deploy",
		Data: event,
	})
}

func (v *Vm) handleBrc20Mint(input *model.Input, output *model.Output, r *model.Content, s *state.State, receipt *model.Receipt) {
	name := strings.ToLower(r.Name)
	if s.IsEmpty(name) {
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
	maximum := s.Get("info:maximum").(int64)
	totalSupply := s.Get("info:total_supply").(int64)
	if totalSupply >= maximum {
		return
	}
	if totalSupply+amount > maximum {
		amount = maximum - totalSupply
	}
	//
	s.Set("info:total_supply", totalSupply+amount)
	balance := s.Get("balance:" + input.Address).(int64)
	s.Set("balance:"+input.Address, balance+amount)
	//
	receipt.Status = 1
	event := make([]string, 3)
	event[0] = input.Address
	event[1] = input.Address
	event[2] = r.Amount
	receipt.Events = append(receipt.Events, model.Event{
		Name: name,
		Id:   "transfer",
		Data: event,
	})
}

func (v *Vm) handleBrc20Transfer(input *model.Input, output *model.Output, r *model.Content, s *state.State, receipt *model.Receipt) {
	name := strings.ToLower(r.Name)
	if s.IsEmpty(name) {
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
	event := make([]string, 3)
	event[0] = fromAddress
	event[1] = toAddress
	event[2] = r.Amount
	receipt.Events = append(receipt.Events, model.Event{
		Name: name,
		Id:   "transfer",
		Data: event,
	})
}
