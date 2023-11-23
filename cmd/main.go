package cmd

import (
	"context"
	"fmt"
	"github.com/motoko9/config"
	"github.com/motoko9/indexer"
	"time"
)

func main() {
	ctx := context.Background()
	indexer := indexer.New(ctx, &config.Config{
		OrdNode: config.OrdNode{
			Url: "http://10.5.20.37:8877",
		},
		BitcoinNode: config.BitcoinNode{
			Url:  "10.5.20.37:8332",
			User: "bitcoinrpc-test",
			Pass: "BA2z8yPRLp2VCMdsdWkUBtqUvvioWasLrHqu88cNK1234-test",
		},
		Postgres: config.Postgres{
			Connect: "host=127.0.0.1 user=postgres password= dbname=ord port=5432 sslmode=disable",
		},
	})
	go func() {
		for i := 0; i < 10000; i++ {
			time.Sleep(time.Second * 100)
			fmt.Printf("syncer height: %d\n", indexer.SyncerHeight())
		}
	}()
	indexer.Service()
}
