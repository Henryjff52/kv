package main

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

// type KVStoreApplication struct{}

var _ abcitypes.Application = (*KVStoreApplication)(nil)

// func NewKVStoreApplication() *KVStoreApplication {
// 	return &KVStoreApplication{}
// }

type KVStoreApplication struct {
	db           *badger.DB
	currentBatch *badger.Txn
}

func NewKVStoreApplication(db *badger.DB) *KVStoreApplication {
	return &KVStoreApplication{
		db: db,
	}
}

func (KVStoreApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

// func (KVStoreApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
// 	return abcitypes.ResponseDeliverTx{Code: 0}
// }

// func (KVStoreApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
// 	return abcitypes.ResponseCheckTx{Code: 0}
// }

func (app *KVStoreApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	fmt.Println("DeliverTx", req)
	return abcitypes.ResponseDeliverTx{Code: 1}

	// code := app.isValid(req.Tx)
	// if code != 0 {
	// 	return abcitypes.ResponseDeliverTx{Code: code}
	// }

	// parts := bytes.Split(req.Tx, []byte("="))
	// key, value := parts[0], parts[1]
	// fmt.Println(string(key), string(value))
	// err := app.currentBatch.Set(key, value)
	// if err != nil {
	// 	panic(err)
	// }

	// return abcitypes.ResponseDeliverTx{Code: 0}
}

// func (KVStoreApplication) Commit() abcitypes.ResponseCommit {
// 	return abcitypes.ResponseCommit{}
// }

func (app *KVStoreApplication) Commit() abcitypes.ResponseCommit {
	fmt.Println("Commit")
	app.currentBatch.Commit()
	return abcitypes.ResponseCommit{Data: []byte{}}
}

// func (KVStoreApplication) Query(req abcitypes.RequestQuery) abcitypes.ResponseQuery {
// 	return abcitypes.ResponseQuery{Code: 0}
// }

func (app *KVStoreApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	fmt.Println("Query", reqQuery)

	resQuery.Key = reqQuery.Data
	err := app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(reqQuery.Data)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == badger.ErrKeyNotFound {
			resQuery.Log = "does not exist"
		} else {
			return item.Value(func(val []byte) error {
				resQuery.Log = "exists"
				resQuery.Value = val
				resQuery.Info = string(resQuery.Key) + "=" + string(resQuery.Value)
				return nil
			})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return
}

func (KVStoreApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

// func (KVStoreApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
// 	return abcitypes.ResponseBeginBlock{}
// }

func (app *KVStoreApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	fmt.Println("BeginBlock")
	time.Sleep(time.Second * 5)

	app.currentBatch = app.db.NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

func (KVStoreApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

func (KVStoreApplication) ListSnapshots(abcitypes.RequestListSnapshots) abcitypes.ResponseListSnapshots {
	return abcitypes.ResponseListSnapshots{}
}

func (KVStoreApplication) OfferSnapshot(abcitypes.RequestOfferSnapshot) abcitypes.ResponseOfferSnapshot {
	return abcitypes.ResponseOfferSnapshot{}
}

func (KVStoreApplication) LoadSnapshotChunk(abcitypes.RequestLoadSnapshotChunk) abcitypes.ResponseLoadSnapshotChunk {
	return abcitypes.ResponseLoadSnapshotChunk{}
}

func (KVStoreApplication) ApplySnapshotChunk(abcitypes.RequestApplySnapshotChunk) abcitypes.ResponseApplySnapshotChunk {
	return abcitypes.ResponseApplySnapshotChunk{}
}
