package indexer

import (
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/hashicorp/go-hclog"
	"github.com/motoko9/config"
	"github.com/motoko9/db"
	"github.com/motoko9/model"
	"github.com/motoko9/ord"
	"github.com/motoko9/state"
	"github.com/motoko9/syncer"
	"github.com/motoko9/vm"
	"golang.org/x/net/context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Indexer struct {
	ctx       context.Context
	log       hclog.Logger
	ordClient *ord.Client
	btcClient *rpcclient.Client
	dao       *db.Dao
	syncer    *syncer.Syncer
	vm        *vm.Vm
}

func New(ctx context.Context, cfg *config.Config) *Indexer {
	ordClient := ord.New(cfg.OrdNode.Url)
	//
	btcClient, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         cfg.BitcoinNode.Url,
		User:         cfg.BitcoinNode.User,
		Pass:         cfg.BitcoinNode.Pass,
		HTTPPostMode: true,
	}, nil)
	if err != nil {
		panic(err)
	}

	//
	// initialize database
	Logger := logger.Default
	if true {
		Logger = Logger.LogMode(logger.Info)
	}
	dbInstance, err := gorm.Open(postgres.Open(cfg.Postgres.Connect), &gorm.Config{Logger: Logger})
	if err != nil {
		panic(err)
	}
	//
	dbInstance.AutoMigrate(&db.Brc20Transaction{}, &db.Brc20Event{}, &db.Brc20Receipt{}, &db.Brc20Info{}, &db.Brc20Balance{}, &db.Inscription{})

	dao := db.NewDao(dbInstance)
	indexer := &Indexer{
		ctx:       ctx,
		log:       hclog.L().Named("indexer"),
		ordClient: ordClient,
		btcClient: btcClient,
		dao:       dao,
	}
	//
	s := syncer.New(ordClient, btcClient, 779832, indexer, indexer.log)
	indexer.syncer = s
	v := vm.New(indexer.log)
	indexer.vm = v
	return indexer
}

func (indexer *Indexer) Service() {
	indexer.Start()
	<-indexer.ctx.Done()
	indexer.Stop()
}

func (indexer *Indexer) Start() {
	indexer.syncer.Start()
	indexer.log.Info("indexer start successful......")
}

func (indexer *Indexer) Stop() {
}

func (indexer *Indexer) OnTransactions(height uint64, txs []*model.Transaction) error {
	indexer.log.Info("OnTransactions", "height", height, "transaction size", len(txs))
	//
	receipts := make([]*model.Receipt, len(txs))
	s := state.New(indexer.dao)
	for j, tx := range txs {
		receipts[j] = indexer.vm.Execute(s, tx)
	}
	//
	s.Commit()
	// save to db
	// todo batch
	brc20Receipts := make([]*db.Brc20Receipt, 0)
	brc20Transactions := make([]*db.Brc20Transaction, 0)
	events := make([]*db.Brc20Event, 0)
	for i, receipt := range receipts {
		if receipt == nil {
			continue
		}
		if receipt.Msg == "not support" {
			continue
		}
		tx := txs[i]
		//
		for _, event := range receipt.Events {
			events = append(events, &db.Brc20Event{
				Brc20:         event.Name,
				InscriptionId: event.Id,
				Data1:         event.Data[0],
				Data2:         event.Data[1],
				Data3:         event.Data[2],
			})
		}
		brc20Receipts = append(brc20Receipts, &db.Brc20Receipt{
			Hash:          receipt.Hash,
			InscriptionId: receipt.InscriptionId,
			Status:        receipt.Status,
			Msg:           receipt.Msg,
			//Events:        events,
		})
		brc20Transactions = append(brc20Transactions, &db.Brc20Transaction{
			Hash:          tx.Hash,
			InscriptionId: tx.InscriptionId,
			ContentLength: tx.Inscription.ContentLength,
			ContentType:   tx.Inscription.ContentType,
			Content:       tx.Inscription.Content,
			Timestamp:     tx.Inscription.Timestamp,
			Height:        height,
		})
	}
	//
	indexer.dao.Brc20Receipt().Save(brc20Receipts)
	indexer.dao.Brc20Transaction().Save(brc20Transactions)
	indexer.dao.Brc20Event().Save(events)

	return nil
}

func (indexer *Indexer) SyncerHeight() uint64 {
	return indexer.syncer.Height()
}
