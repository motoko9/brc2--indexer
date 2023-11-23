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

func (v *Vm) Execute(s *state.State, transaction *model.Transaction) *model.Context {
	if transaction.Inscription.ContentType != "text/plain;charset=utf-8" {
		return nil
	}
	//
	var r model.Content
	err := json.Unmarshal(transaction.Inscription.Content, &r)
	if err != nil {
		return nil
	}
	//
	if r.Proto != "brc-20" {
		return nil
	}
	//
	if r.Name != "ordi" {
		return nil
	}
	//
	c := &model.Context{
		Output:      transaction.Output,
		Inscription: transaction.Inscription,
		Content:     r,
		Status:      0,
		Msg:         "",
	}
	v.ExecuteBrc20(c, s)
	return c
}

func (v *Vm) ExecuteBrc20(c *model.Context, s *state.State) {
	r := c.Content
	switch r.Operation {
	case "deploy":
		v.handleBrc20Deploy(c, s)
	case "mint":
		v.handleBrc20Mint(c, s)
	case "transfer":
		v.handleBrc20Transfer(c, s)
	default:
		c.Status = 0
		c.Msg = InvalidFunction
		return
	}
}

func (v *Vm) handleBrc20Deploy(c *model.Context, s *state.State) {
	inscription := c.Inscription
	if s.HasInscription(inscription.InscriptionId) {
		c.Status = 0
		c.Msg = DuplicateResource
		return
	}
	output := c.Output
	s.CreateInscription(inscription.InscriptionId, inscription.ContentLength, inscription.ContentType, inscription.Content, output.Address)
	//
	r := c.Content
	name := strings.ToLower(r.Name)
	if s.HasBr20(name) {
		c.Status = 0
		c.Msg = DuplicateResource
		return
	}
	//
	var err error
	if r.Maximum == "" {
		c.Status = 0
		c.Msg = InvalidParameter
		return
	}
	maximum, err := strconv.ParseInt(r.Maximum, 10, 64)
	if err != nil {
		c.Status = 0
		c.Msg = InvalidParameter
		return
	}

	decimal := int64(1)
	if r.Decimal != "" {
		decimal, err = strconv.ParseInt(r.Decimal, 10, 64)
		if err != nil {
			c.Status = 0
			c.Msg = InvalidParameter
			return
		}
	}
	limit := int64(math.MaxInt64)
	if r.Limit != "" {
		limit, err = strconv.ParseInt(r.Limit, 10, 64)
		if err != nil {
			c.Status = 0
			c.Msg = InvalidParameter
			return
		}
	}
	//
	s.CreateBrc20(r.Name)
	s.SetCurrentBrc20(r.Name)
	s.Set("info:name", name)
	s.Set("info:decimal", decimal)
	s.Set("info:maximum", maximum)
	s.Set("info:limit", limit)
	s.Set("info:total_supply", int64(0))
	//
	c.Status = 1
	event := make([]string, 3)
	event[0] = r.Decimal
	event[1] = r.Maximum
	event[2] = r.Limit
	c.Event = model.Event{
		Name: name,
		Id:   "deploy",
		Data: event,
	}
}

func (v *Vm) handleBrc20Mint(c *model.Context, s *state.State) {
	inscription := c.Inscription
	if s.HasInscription(inscription.InscriptionId) {
		c.Status = 0
		c.Msg = DuplicateResource
		return
	}
	//
	output := c.Output
	s.CreateInscription(inscription.InscriptionId, inscription.ContentLength, inscription.ContentType, inscription.Content, output.Address)
	//
	r := c.Content
	name := strings.ToLower(r.Name)
	if !s.HasBr20(name) {
		c.Status = 0
		c.Msg = UnknownResource
		return
	}
	//
	s.SetCurrentBrc20(r.Name)
	//
	amount, err := strconv.ParseInt(r.Amount, 10, 64)
	if err != nil {
		c.Status = 0
		c.Msg = InvalidParameter
		return
	}
	limit := s.Get("info:limit").(int64)
	if amount > limit {
		c.Status = 0
		c.Msg = InvalidParameter
		return
	}
	maximum := s.Get("info:maximum").(int64)
	totalSupply := s.Get("info:total_supply").(int64)
	if totalSupply >= maximum {
		c.Status = 0
		c.Msg = InvalidParameter
		return
	}
	if totalSupply+amount > maximum {
		amount = maximum - totalSupply
	}
	//
	s.Set("info:total_supply", totalSupply+amount)
	balance := s.Get("balance:" + output.Address).(int64)
	s.Set("balance:"+output.Address, balance+amount)
	//
	c.Status = 1
	event := make([]string, 3)
	event[0] = output.Address
	event[1] = output.Address
	event[2] = r.Amount
	c.Event = model.Event{
		Name: name,
		Id:   "transfer",
		Data: event,
	}
}

func (v *Vm) handleBrc20Transfer(c *model.Context, s *state.State) {
	inscription := c.Inscription
	if !s.HasInscription(inscription.InscriptionId) {
		v.handleBrc20TransferStep1(c, s)
		return
	} else {
		v.handleBrc20TransferStep2(c, s)
		return
	}
}

func (v *Vm) handleBrc20TransferStep1(c *model.Context, s *state.State) {
	//
	inscription := c.Inscription
	output := c.Output
	s.CreateInscription(inscription.InscriptionId, inscription.ContentLength, inscription.ContentType, inscription.Content, output.Address)
	//
	r := c.Content
	name := strings.ToLower(r.Name)
	if !s.HasBr20(name) {
		c.Status = 0
		c.Msg = UnknownResource
		return
	}
	//
	c.Status = 1
}

func (v *Vm) handleBrc20TransferStep2(c *model.Context, s *state.State) {
	//
	inscription := c.Inscription
	number := s.IncreaseInscriptionNumber(inscription.InscriptionId)
	output := c.Output
	oldOwner := s.SetNewOwner(inscription.InscriptionId, output.Address)
	if number != 1 {
		c.Status = 0
		c.Msg = DuplicateResource
		return
	}
	//
	r := c.Content
	name := strings.ToLower(r.Name)
	if !s.HasBr20(name) {
		c.Status = 0
		c.Msg = UnknownResource
		return
	}
	//
	fromAddress := oldOwner
	toAddress := output.Address
	amount, err := strconv.ParseInt(r.Amount, 10, 64)
	if err != nil {
		c.Status = 0
		c.Msg = InvalidParameter
		return
	}
	fromBalance := s.Get("balance:" + fromAddress).(int64)
	toBalance := s.Get("balance:" + toAddress).(int64)
	if fromBalance < amount {
		c.Status = 0
		c.Msg = InvalidParameter
		return
	}
	s.Set("balance:"+fromAddress, fromBalance-amount)
	s.Set("balance:"+toAddress, toBalance+amount)
	//
	c.Status = 1
	event := make([]string, 3)
	event[0] = fromAddress
	event[1] = toAddress
	event[2] = r.Amount
	c.Event = model.Event{
		Name: name,
		Id:   "transfer",
		Data: event,
	}
}
