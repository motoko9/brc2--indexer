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
	indexer := New(ctx, &config.Config{OrdNode: config.OrdNode{Url: "https://ordinals.com"}})
	indexer.Start()
	//
	for i := 0; i < 10000; i++ {
		time.Sleep(time.Second * 100)
		fmt.Printf("syncer height: %d\n", indexer.SyncerHeight())
	}
}
