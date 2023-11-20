package syncer

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/hashicorp/go-hclog"
	"github.com/motoko9/model"
	"github.com/motoko9/ord"
	"strings"
	"time"
)

type Callback interface {
	OnTransactions(height uint64, txs []*model.Transaction) error
}

type Syncer struct {
	ordClient *ord.Client
	btcClient *rpcclient.Client
	height    uint64
	cb        Callback
	log       hclog.Logger
}

func New(ordClient *ord.Client, btcClient *rpcclient.Client, height uint64, cb Callback, log hclog.Logger) *Syncer {
	i := &Syncer{
		ordClient: ordClient,
		btcClient: btcClient,
		height:    height,
		cb:        cb,
		log:       log.Named("syncer"),
	}
	return i
}

func (syncer *Syncer) Start() {
	go syncer.process()
}

func (syncer *Syncer) sync() bool {
	latestHeight, err := syncer.ordClient.BlockHeight()
	if err != nil {
		syncer.log.Error("ordClient.BlockHeight", "error", err)
		return false
	}
	syncer.log.Info("sync", "latest height", latestHeight, "sync height", syncer.height)
	for syncer.height < latestHeight {
		syncer.log.Info("sync", "height", syncer.height)
		//
		txs := make([]*model.Transaction, 0)
		ids, err := syncer.ordClient.InscriptionsByBlock(syncer.height)
		if err != nil {
			syncer.log.Error("ordClient.InscriptionsByBlock", "error", err)
			return false
		}
		//
		for _, id := range ids {
			s, err := syncer.ordClient.InscriptionById(id)
			if err != nil {
				syncer.log.Error("ordClient.InscriptionById", "error", err)
				return false
			}
			c, err := syncer.ordClient.InscriptionContent(id)
			if err != nil {
				syncer.log.Error("ordClient.InscriptionContent", "error", err)
				return false
			}
			txid := s.InscriptionId[0:64]
			n := MustUint64(s.InscriptionId[65:])
			//
			tx, err := syncer.ordClient.Tx(txid)
			if err != nil {
				syncer.log.Error("ordClient.Tx", "error", err)
				return false
			}
			// todo
			items := strings.Split(tx.Inputs[0].Id, ":")
			inputTxHash := items[0]
			inputTxN := MustUint64(items[1])
			inputTx, err := syncer.ordClient.Tx(inputTxHash)
			if err != nil {
				syncer.log.Error("ordClient.Tx2", "error", err)
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
		if syncer.cb != nil {
			if err := syncer.cb.OnTransactions(syncer.height, txs); err != nil {
				syncer.log.Error("cb.OnTransactions", "error", err)
				return false
			}
		}
		//
		syncer.height++
	}
	return true
}

func (syncer *Syncer) process() {
	process := func() (exit bool) {
		defer func() {
			if r := recover(); r != nil {
				exit = false
			}
		}()
		syncer.log.Info("syncer process start......")
		ticker := time.NewTicker(time.Second * 1)
		for {
			select {
			case <-ticker.C:
				if !syncer.sync() {
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

func (syncer *Syncer) Height() uint64 {
	return syncer.height
}
