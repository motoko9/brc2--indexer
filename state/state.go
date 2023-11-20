package state

import (
	"github.com/motoko9/db"
	"github.com/motoko9/state/brc20"
)

type State struct {
	dao   *db.Dao
	name  string
	brc20 *brc20.State
}

func New(dao *db.Dao) *State {
	s := &State{
		dao: dao,
	}
	return s
}

func (s *State) Reload(name string) {
	s.name = name
	s.brc20 = brc20.Load(s.dao, name)
}

func (s *State) IsEmpty() bool {
	return s.brc20 == nil
}

func (s *State) Create(name string) {
	s.brc20 = brc20.New(s.dao, name)
}

func (s *State) Set(key string, value interface{}) {
	s.brc20.Set(key, value)
}

func (s *State) Get(key string) interface{} {
	return s.brc20.Get(key)
}

func (s *State) Commit() {

}
