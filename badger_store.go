package i18n_manager

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateBadgerStore(rootPath string) (Store, error) {
	if file, err := os.Stat(rootPath); (nil != err && os.IsNotExist(err)) || (nil != file && !file.IsDir()) {
		if err = os.MkdirAll(rootPath, os.FileMode(0660)); nil != err {
			panic(err)
		}
	}

	db, err := openDB(rootPath)
	if err != nil {
		return nil, err
	}

	return &store{db}, nil
}

func openDB(rootPath string) (*badger.DB, error) {
	opts := badger.DefaultOptions
	opts.Dir = rootPath
	opts.ValueDir = rootPath
	db, err := badger.Open(opts)
	if err != nil {
		if strings.Contains(err.Error(), "LOCK") {
			log.Println("database locked, probably due to improper shutdown")
			if db, err := retryOpen(rootPath, opts); err == nil {
				log.Println("database unlocked, value log truncated")
				return db, nil
			}
			log.Println("could not unlock database:", err)

		}
		return nil, err
	}
	return db, nil
}

func retryOpen(dir string, originalOpts badger.Options) (*badger.DB, error) {
	lockPath := filepath.Join(dir, "LOCK")
	if err := os.Remove(lockPath); err != nil {
		return nil, fmt.Errorf(`removing "LOCK": %s`, err)
	}
	retryOpts := originalOpts
	retryOpts.Truncate = true
	db, err := badger.Open(retryOpts)
	return db, err
}

type store struct {
	db *badger.DB
}

func (s *store) Save(items []*I18NItem) error {
	return s.db.Update(func(txn *badger.Txn) error {
		for _, item := range items {
			key := []byte(item.Key)
			if bytes, err := json.Marshal(item); nil != err {
				return err
			} else if err = txn.Set(key, bytes); nil != err {
				return err
			}
		}

		return nil
	})
}

func (s *store) LoadAll() (items []*I18NItem, err error) {
	items = make([]*I18NItem, 0, 100)
	err = s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte{0}
		for it.Seek(prefix); ; it.Next() {
			if !it.Valid() {
				break
			}
			bytes, err := it.Item().Value()
			if nil == err {
				item := &I18NItem{}
				if err = json.Unmarshal(bytes, item); nil == err {
					items = append(items, item)
				}
			}
		}

		return nil
	})

	return
}

func (s *store) Shutdown() error {
	return s.db.Close()
}
