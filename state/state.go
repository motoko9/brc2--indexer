package state

import (
	"github.com/motoko9/db"
	"github.com/motoko9/state/brc20"
)

type State struct {
	dao          *db.Dao
	latestHeight int64
	current      string
	brc20        map[string]*brc20.State
	inscriptions map[string]*db.Inscription
}

func New(dao *db.Dao) *State {
	s := &State{
		dao:          dao,
		brc20:        make(map[string]*brc20.State),
		inscriptions: make(map[string]*db.Inscription),
	}
	return s
}

func (s *State) HasInscription(inscriptionId string) bool {
	inscription, ok := s.inscriptions[inscriptionId]
	if ok {
		return true
	}
	inscription, err := s.dao.Inscription().Find(inscriptionId)
	if err != nil {
		return false
	}
	s.inscriptions[inscriptionId] = inscription
	return true
}

func (s *State) CreateInscription(inscriptionId string, contentLength uint64, contentType string, content []byte, owner string) {
	s.inscriptions[inscriptionId] = &db.Inscription{
		InscriptionId:     "",
		ContentLength:     contentLength,
		ContentType:       contentType,
		Content:           content,
		Owner:             owner,
		InscriptionNumber: 0,
	}
}

func (s *State) IncreaseInscriptionNumber(inscriptionId string) int64 {
	inscription, ok := s.inscriptions[inscriptionId]
	if !ok {
		return 0
	}
	inscription.InscriptionNumber += 1
	return inscription.InscriptionNumber
}

func (s *State) SetNewOwner(inscriptionId string, address string) string {
	inscription, ok := s.inscriptions[inscriptionId]
	if !ok {
		return ""
	}
	oldOwner := inscription.Owner
	inscription.Owner = address
	return oldOwner
}

func (s *State) HasBr20(name string) bool {
	_, ok := s.brc20[name]
	if ok {
		return true
	}
	s.brc20[name] = brc20.Load(s.dao, name)
	item, ok := s.brc20[name]
	return ok && item != nil
}

func (s *State) CreateBrc20(name string) {
	item, ok := s.brc20[name]
	if ok && item != nil {
		return
	}
	s.brc20[name] = brc20.New(s.dao, name)
}

func (s *State) SetCurrentBrc20(name string) {
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
	inscriptions := make([]*db.Inscription, 0)
	for _, item := range s.inscriptions {
		inscriptions = append(inscriptions, item)
	}
	s.dao.Inscription().Save(inscriptions)
	s.dao.Sync().UpdateCommitHeight(s.latestHeight)
}

func (s *State) UpdateHeight(height int64) {
	s.latestHeight = height
}
