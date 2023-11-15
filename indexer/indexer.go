package indexer

import (
	"github.com/motoko9/config"
	"github.com/motoko9/model"
	"github.com/motoko9/ord"
	"github.com/motoko9/syncer"
	"github.com/motoko9/vm"
	"golang.org/x/net/context"
)

type Indexer struct {
	ctx       context.Context
	ordClient *ord.Client
	syncer    *syncer.Syncer
	vm        *vm.Vm
}

func New(ctx context.Context, cfg *config.Config) *Indexer {
	ordClient := ord.New(cfg.OrdNode.Url)
	i := &Indexer{
		ctx:       ctx,
		ordClient: ordClient,
	}
	//
	s := syncer.New(ordClient, 779630, i)
	i.syncer = s
	v := vm.New()
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

func (i *Indexer) OnBrc20Transactions(height uint64, txs []*model.Brc20Transaction) error {
	for _, tx := range txs {
		i.vm.Execute(tx)
	}
	return nil
}

func (i *Indexer) SyncerHeight() uint64 {
	return i.syncer.Height()
}
