package indexer

import (
	"github.com/btcsuite/btcd/rpcclient"
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
	ordClient *ord.Client
	btcClient *rpcclient.Client
	dao       *db.Dao
	syncer    *syncer.Syncer
	s         *state.State
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
		return nil
	}

	//
	// initialize database
	Logger := logger.Default
	if true {
		Logger = Logger.LogMode(logger.Info)
	}
	dbInstance, err := gorm.Open(postgres.Open(cfg.Postgres.Connect), &gorm.Config{Logger: Logger})
	if err != nil {
		return nil
	}
	dao := db.NewDao(dbInstance)
	i := &Indexer{
		ctx:       ctx,
		ordClient: ordClient,
		btcClient: btcClient,
		dao:       dao,
	}
	//
	s := syncer.New(ordClient, btcClient, 779630, i)
	i.syncer = s
	st := state.New(dao)
	i.s = st
	v := vm.New(st)
	i.vm = v
	return i
}

func (i *Indexer) Service() {
	i.Start()
	<-i.ctx.Done()
	i.Stop()
}

func (i *Indexer) Start() {
	i.syncer.Start()
}

func (i *Indexer) Stop() {
}

func (i *Indexer) OnBrc20Transactions(height uint64, txs []*model.Transaction) error {
	receipt := make([]*model.Receipt, len(txs))
	for j, tx := range txs {
		receipt[j] = i.vm.Execute(i.s, tx)
	}
	//
	i.s.Commit()

	// save to db
	// todo batch

	return nil
}

func (i *Indexer) SyncerHeight() uint64 {
	return i.syncer.Height()
}
