package brc20

import (
	"github.com/motoko9/db"
	"strings"
)

type State struct {
	dao      *db.Dao
	info     *db.Brc20Info
	balances map[string]int64
}

func Load(dao *db.Dao, name string) *State {
	s := &State{
		dao: dao,
	}
	//
	b := s.dao.Brc20Info()
	info, err := b.Find(name)
	if err != nil {
		return nil
	}
	if info == nil {
		return nil
	}
	s.info = info
	s.balances = make(map[string]int64)
	return s
}

func New(dao *db.Dao, name string) *State {
	s := &State{
		dao: dao,
	}
	s.info = &db.Brc20Info{
		Name: name,
	}
	s.balances = make(map[string]int64)
	return s
}

func (s *State) Get(key string) interface{} {
	items := strings.Split(key, ":")
	switch items[0] {
	case "info":
		switch items[1] {
		case "name":
			return s.info.Name
		case "decimal":
			return s.info.Decimal
		case "maximum":
			return s.info.Maximum
		case "limit":
			return s.info.Limit
		case "total_supply":
			return s.info.TotalSupply
		}
	case "balance":
		return s.balances[items[1]]
	}
	return nil
}

func (s *State) Set(key string, value interface{}) {
	items := strings.Split(key, ":")
	switch items[0] {
	case "info":
		switch items[1] {
		case "name":
			s.info.Name = value.(string)
		case "decimal":
			s.info.Decimal = value.(int64)
		case "maximum":
			s.info.Maximum = value.(int64)
		case "limit":
			s.info.Limit = value.(int64)
		case "total_supply":
			s.info.TotalSupply = value.(int64)
		}
	case "balance":
		s.balances[items[1]] = value.(int64)
	}
}

func (s *State) Commit() {
	s.dao.Brc20Info().Save(s.info)
	balances := make([]*db.Brc20Balance, 0)
	for address, balance := range s.balances {
		balances = append(balances, &db.Brc20Balance{
			Name:    s.info.Name,
			Address: address,
			Balance: balance,
		})
	}
	s.dao.Brc20Balance().Save(balances)
}
