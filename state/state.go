package state

import (
	"github.com/motoko9/db"
	"github.com/motoko9/state/brc20"
)

type State struct {
	dao     *db.Dao
	current string
	brc20   map[string]*brc20.State
}

func New(dao *db.Dao) *State {
	s := &State{
		dao: dao,
	}
	return s
}

func (s *State) Load(name string) {
	_, ok := s.brc20[name]
	if ok {
		s.current = name
		return
	}
	s.brc20[name] = brc20.Load(s.dao, name)
	s.current = name
}

func (s *State) IsEmpty(name string) bool {
	_, ok := s.brc20[name]
	return !ok
}

func (s *State) Create(name string) {
	_, ok := s.brc20[name]
	if ok {
		s.current = name
		return
	}
	s.brc20[name] = brc20.New(s.dao, name)
	s.current = name
}

func (s *State) Set(key string, value interface{}) {
	brc20, ok := s.brc20[s.current]
	if !ok {
		return
	}
	brc20.Set(key, value)
}

func (s *State) Get(key string) interface{} {
	brc20, ok := s.brc20[s.current]
	if !ok {
		return nil
	}
	return brc20.Get(key)
}

func (s *State) Commit() {
	for _, brc20 := range s.brc20 {
		brc20.Commit()
	}
}
