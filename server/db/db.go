package db

import (
	"errors"
	"github.com/dgraph-io/badger/v3"
	"log"
	"server/db/keys"
	"server/db/values"
	"server/env"
)

var db *badger.DB

func init() {
	var err error
	if db, err = badger.Open(badger.DefaultOptions(env.DatabaseDirectory())); err != nil {
		panic(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(keys.PostsCount()); !errors.Is(err, badger.ErrKeyNotFound) {
			return err
		}
		if err := txn.Set(keys.PostsCount(), values.PostsCount(0)); err != nil {
			return err
		}
		log.Println("[db] created field PostsCount")
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func Close() {
	if err := db.Close(); err != nil {
		log.Println("[db] close error:", err)
	}
}
