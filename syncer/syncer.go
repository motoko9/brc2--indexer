package syncer

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/motoko9/model"
	"github.com/motoko9/ord"
	"strings"
	"time"
)

type Callback interface {
	OnBrc20Transactions(height uint64, txs []*model.Transaction) error
}

type Syncer struct {
	ordClient *ord.Client
	btcClient *rpcclient.Client
	height    uint64
	cb        Callback
}

func New(ordClient *ord.Client, btcClient *rpcclient.Client, height uint64, cb Callback) *Syncer {
	i := &Syncer{
		ordClient: ordClient,
		btcClient: btcClient,
		height:    height,
		cb:        cb,
	}
	return i
}

func (i *Syncer) Start() {
	go i.process()
}

func (i *Syncer) sync() bool {
	latestHeight, err := i.ordClient.BlockHeight()
	if err != nil {
		return false
	}
	for i.height < latestHeight {
		txs := make([]*model.Transaction, 0)
		ids, err := i.ordClient.InscriptionsByBlock(i.height)
		if err != nil {
			return false
		}
		//
		for _, id := range ids {
			s, err := i.ordClient.InscriptionById(id)
			if err != nil {
				return false
			}
			c, err := i.ordClient.InscriptionContent(id)
			if err != nil {
				return false
			}
			txid := s.InscriptionId[0:64]
			n := MustUint64(s.InscriptionId[65:])
			//
			tx, err := i.ordClient.Tx(txid)
			if err != nil {
				return false
			}
			// todo
			items := strings.Split(tx.Inputs[0].Id, ":")
			inputTxHash := items[0]
			inputTxN := MustUint64(items[1])
			inputTx, err := i.ordClient.Tx(inputTxHash)
			if err != nil {
				return false
			}
			//
			brc20Tx := model.Transaction{
				Hash: txid,
				Input: model.Input{
					Hash:    inputTxHash,
					N:       inputTxN,
					Address: inputTx.Outputs[inputTxN].Address,
					Value:   MustUint64(inputTx.Outputs[inputTxN].Value),
				},
				Output: model.Output{
					Value:   MustUint64(tx.Outputs[n].Value),
					N:       n,
					Address: tx.Outputs[n].Address,
				},
				Inscription: model.Inscription{
					Address:           s.Address,
					ContentLength:     s.ContentLength,
					ContentType:       s.ContentType,
					Content:           c,
					GenesisFee:        s.GenesisFee,
					GenesisHeight:     s.GenesisHeight,
					InscriptionId:     s.InscriptionId,
					InscriptionNumber: s.InscriptionNumber,
					OutputValue:       s.OutputValue,
					SatPoint:          s.SatPoint,
					Timestamp:         s.Timestamp,
				},
			}
			txs = append(txs, &brc20Tx)
		}
		//
		if i.cb != nil {
			if err := i.cb.OnBrc20Transactions(i.height, txs); err != nil {
				return false
			}
		}
		//
		i.height++
	}
	return true
}

func (i *Syncer) process() {
	process := func() (exit bool) {
		defer func() {
			if r := recover(); r != nil {
				exit = false
			}
		}()
		ticker := time.NewTicker(time.Second * 1)
		for {
			select {
			case <-ticker.C:
				if !i.sync() {
					return
				}
			}
		}
	}
	for {
		if exit := process(); exit {
			return
		}
		time.Sleep(time.Second * 5)
	}
}

func (s *Syncer) Height() uint64 {
	return s.height
}
