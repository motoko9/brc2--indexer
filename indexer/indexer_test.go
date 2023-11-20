package indexer

import (
	"fmt"
	"github.com/motoko9/config"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestIndexer_Start(t *testing.T) {
	ctx := context.Background()
	indexer := New(ctx, &config.Config{
		OrdNode: config.OrdNode{
			Url: "http://10.5.20.37:8877",
		},
		BitcoinNode: config.BitcoinNode{
			Url:  "http://10.5.20.37:8877",
			User: "bitcoinrpc-test",
			Pass: "BA2z8yPRLp2VCMdsdWkUBtqUvvioWasLrHqu88cNK1234-test",
		},
		Postgres: config.Postgres{
			Connect: "host=127.0.0.1 user=tangaoyuan password=123456 dbname=wallet port=5432 sslmode=disable",
		},
	})
	indexer.Start()
	//
	for i := 0; i < 10000; i++ {
		time.Sleep(time.Second * 100)
		fmt.Printf("syncer height: %d\n", indexer.SyncerHeight())
	}
}
