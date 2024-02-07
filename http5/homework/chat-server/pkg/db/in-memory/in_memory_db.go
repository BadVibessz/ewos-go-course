package in_memory

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

var (
	ErrNotExistedRow   = errors.New("no such row")
	ErrNotExistedTable = errors.New("no such table")
)

const (
	writePerm = 0o664
)

type Table = *orderedmap.OrderedMap[string, any]

type InMemoryDB interface {
	CreateTable(name string)
	GetTable(name string) (Table, error)
	DropTable(name string)

	AddRow(table string, identifier string, row any) error
	AlterRow(table string, identifier string, newRow any) error
	DropRow(table string, identifier string) error
	GetRow(table string, identifier string) (any, error)
	GetAllRows(table string, offset, limit int) ([]any, error)

	GetRowsCount(table string) (int, error)

	Clear()
}

type InMemDB struct {
	Tables map[string]Table
	m      *sync.RWMutex
}

func NewInMemDB(ctx context.Context, savePath string) (*InMemDB, <-chan any) {
	db := InMemDB{
		Tables: make(map[string]Table),
		m:      &sync.RWMutex{},
	}

	savedChan := make(chan any, 1)

	go func() {
		<-ctx.Done()
		db.Save(savePath, savedChan)
	}()

	return &db, savedChan
}

func NewInMemDBFromJSON(ctx context.Context, jsonState string, savePath string) (*InMemDB, <-chan any, error) {
	tables := make(map[string]Table)

	err := json.Unmarshal([]byte(jsonState), &tables) // todo: not unmarshalls embedded map
	if err != nil {
		return nil, nil, err
	}

	db := InMemDB{
		Tables: tables,
		m:      &sync.RWMutex{},
	}

	savedChan := make(chan any)

	go func() {
		<-ctx.Done()
		db.Save(savePath, savedChan)
	}()

	return &db, savedChan, nil
}

func (db *InMemDB) Save(path string, doneChan chan any) {
	bytes, err := json.Marshal(db.Tables)
	if err != nil {
		doneChan <- err
		return // todo: log?
	}

	err = os.WriteFile(path, bytes, writePerm)
	if err != nil {
		doneChan <- err
		return
	}

	doneChan <- "ok"
}

func (db *InMemDB) CreateTable(name string) {
	db.m.Lock()
	defer db.m.Unlock()

	db.Tables[name] = orderedmap.New[string, any]()
}

func (db *InMemDB) GetTable(name string) (Table, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	t, ok := db.Tables[name]
	if ok {
		return t, nil
	}

	return nil, ErrNotExistedTable
}

func (db *InMemDB) DropTable(name string) {
	db.m.Lock()
	defer db.m.Unlock()

	delete(db.Tables, name)
}

func (db *InMemDB) Clear() {
	db.m.Lock()
	defer db.m.Unlock()

	db.Tables = make(map[string]Table)
}

func (db *InMemDB) AddRow(table string, identifier string, row any) error {
	t, err := db.GetTable(table)
	if err != nil {
		return err
	}

	db.m.Lock()
	defer db.m.Unlock()

	t.Set(identifier, row)

	return nil
}

func (db *InMemDB) AlterRow(table string, identifier string, newRow any) error {
	db.m.Lock()
	defer db.m.Unlock()

	t, err := db.GetTable(table)
	if err != nil {
		return err
	}

	_, existed := t.Get(identifier)
	if !existed {
		return ErrNotExistedRow
	}

	t.Set(identifier, newRow) // todo: test if it's replaces existing value

	return nil
}

func (db *InMemDB) GetRow(table string, identifier string) (any, error) {
	t, err := db.GetTable(table)
	if err != nil {
		return 0, err
	}

	db.m.RLock()
	defer db.m.RUnlock()

	row, exist := t.Get(identifier)
	if !exist {
		return nil, ErrNotExistedRow
	}

	return row, nil
}

func (db *InMemDB) GetAllRows(table string, offset, limit int) ([]any, error) {
	t, err := db.GetTable(table)
	if err != nil {
		return nil, err
	}

	res := make([]any, 0, t.Len())

	db.m.RLock()
	defer db.m.RUnlock()

	count := 0

	// iterating pairs from oldest to newest:
	for pair := t.Oldest(); pair != nil; pair = pair.Next() {
		if count >= offset {
			res = append(res, pair.Value)
		}

		if len(res) == limit {
			break
		}

		count++
	}

	return res, nil
}

func (db *InMemDB) GetRowsCount(table string) (int, error) {
	t, err := db.GetTable(table)
	if err != nil {
		return 0, err
	}

	db.m.RLock()
	defer db.m.RUnlock()

	return t.Len(), nil
}

func (db *InMemDB) DropRow(table string, identifier string) error {
	t, err := db.GetTable(table)
	if err != nil {
		return err
	}

	db.m.Lock()
	defer db.m.Unlock()

	t.Delete(identifier)

	return nil
}
