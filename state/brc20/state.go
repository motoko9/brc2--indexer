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
		Tick: name,
	}
	s.balances = make(map[string]int64)
	return s
}

func (s *State) Get(key string) interface{} {
	items := strings.Split(key, ":")
	switch items[0] {
	case "info":
		switch items[1] {
		case "decimal":
			return s.info.Decimal
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
		case "decimal":
			s.info.Decimal = value.(int64)
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
			Tick:    s.info.Tick,
			Address: address,
			Balance: balance,
		})
	}
	s.dao.Brc20Balance().Save(balances)
}
