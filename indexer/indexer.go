package indexer

import (
	"github.com/braintree/manners"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/gin-gonic/gin"
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
	"gorm.io/gorm/schema"
	"net"
	"net/http"
	"strconv"
)

type Indexer struct {
	ctx       context.Context
	log       hclog.Logger
	ordClient *ord.Client
	btcClient *rpcclient.Client
	dao       *db.Dao
	engine    *gin.Engine
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
		DisableTLS:   true,
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
	dbInstance, err := gorm.Open(postgres.Open(cfg.Postgres.Connect),
		&gorm.Config{
			Logger: Logger,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "t_",
			},
		})
	if err != nil {
		panic(err)
	}
	//
	err = dbInstance.AutoMigrate(&db.Sync{}, &db.Inscription{}, &db.Brc20Info{}, &db.Brc20Balance{}, &db.Transaction{}, &db.Receipt{})
	if err != nil {
		panic(err)
	}
	dao := db.NewDao(dbInstance)

	//
	indexer := &Indexer{
		ctx:       ctx,
		log:       hclog.L().Named("indexer"),
		ordClient: ordClient,
		btcClient: btcClient,
		engine:    gin.New(),
		dao:       dao,
	}
	//
	sync, err := indexer.dao.Sync().Find()
	if err != nil {
		sync := &db.Sync{
			SyncHeight:   779832,
			CommitHeight: 779832,
		}
		indexer.dao.Sync().Save(sync)
	}
	s := syncer.New(ordClient, btcClient, sync.SyncHeight, indexer, indexer.log)
	indexer.syncer = s
	v := vm.New(indexer.log)
	indexer.vm = v
	return indexer
}

func (indexer *Indexer) Service() {
	indexer.Start()
	//
	rootPath := indexer.engine.Group("/api")
	indexer.APIRoutes(rootPath)

	// setup listener
	ln, err := net.Listen("tcp", ":8089")
	if err != nil {
		return
	}
	manners.Serve(ln, indexer.engine)
	//
	indexer.Stop()
}

func (indexer *Indexer) Start() {
	indexer.syncer.Start()
	indexer.log.Info("indexer start successful......")
}

func (indexer *Indexer) Stop() {
}

func (indexer *Indexer) OnTransactions(height int64, txs []*model.Transaction) error {
	indexer.log.Info("OnTransactions", "height", height, "transaction size", len(txs))
	//
	contexts := make([]*model.Context, len(txs))
	s := state.New(indexer.dao)
	for i, tx := range txs {
		contexts[i] = indexer.vm.Execute(s, tx)
	}
	// save transactions
	inscriptionTransactions := make([]*db.Transaction, 0)
	receipts := make([]*db.Receipt, 0)
	for i, item := range contexts {
		if item == nil {
			continue
		}
		if item.Msg == "not support" {
			continue
		}
		inscriptionTransactions = append(inscriptionTransactions, &db.Transaction{
			Hash:          txs[i].Output.Hash,
			N:             int(txs[i].Output.N),
			InscriptionId: txs[i].Inscription.InscriptionId,
			Timestamp:     height,
			Height:        height,
		})
		receipts = append(receipts, &db.Receipt{
			Hash:          txs[i].Output.Hash,
			N:             int(txs[i].Output.N),
			InscriptionId: txs[i].Inscription.InscriptionId,
			Status:        item.Status,
			Msg:           item.Msg,
			Data1:         item.Event.Name,
			Data2:         item.Event.Id,
			Data3:         item.Event.Data[0],
			Data4:         item.Event.Data[1],
			Data5:         item.Event.Data[2],
		})
	}
	indexer.dao.InscriptionTransaction().Save(inscriptionTransactions)
	indexer.dao.Receipt().Save(receipts)
	indexer.dao.Sync().UpdateSyncHeight(height)
	//
	s.UpdateHeight(height)
	s.Commit()
	return nil
}

func (indexer *Indexer) SyncerHeight() int64 {
	return indexer.syncer.Height()
}

//

func (indexer *Indexer) Height(c *gin.Context) {
	sync, err := indexer.dao.Sync().Find()
	if err != nil {
		_ = c.Error(err)
		return
	} else {
		c.JSON(http.StatusOK, struct {
			SyncHeight   int64 `json:"sync_height"`
			CommitHeight int64 `json:"commit_height"`
		}{
			SyncHeight:   sync.SyncHeight,
			CommitHeight: sync.CommitHeight,
		})
	}
}

func (indexer *Indexer) Receipts(c *gin.Context) {
	heightStr := c.Param("height")
	height, err := strconv.ParseInt(heightStr, 10, 64)
	if err != nil {
		_ = c.Error(err)
		return
	}
	transactions, err := indexer.dao.InscriptionTransaction().FindByHeight(height)
	if err != nil {
		_ = c.Error(err)
		return
	}
	type receipt struct {
		Status int       `json:"status"`
		Msg    string    `json:"msg"`
		Data   [5]string `json:"data"`
	}
	type response struct {
		Hash          string    `json:"hash"`
		InscriptionId string    `json:"inscription_id"`
		Timestamp     int64     `json:"timestamp"`
		Height        int64     `json:"height"`
		Receipts      []receipt `json:"receipts"`
	}
	rsp := make([]*response, 0)
	for _, transaction := range transactions {
		receipts := make([]receipt, 0)
		for _, item := range transaction.Receipts {
			receipts = append(receipts, receipt{
				Status: item.Status,
				Msg:    item.Msg,
				Data:   [5]string{item.Data1, item.Data2, item.Data3, item.Data4, item.Data5},
			})
		}
		rsp = append(rsp, &response{
			Hash:          transaction.Hash,
			InscriptionId: transaction.InscriptionId,
			Timestamp:     transaction.Timestamp,
			Height:        transaction.Height,
			Receipts:      receipts,
		})
	}
	c.JSON(http.StatusOK, rsp)
}

func (indexer *Indexer) APIRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	v1.GET("/height", indexer.Height)
	v1.GET("/receipts/:height", indexer.Receipts)
}
