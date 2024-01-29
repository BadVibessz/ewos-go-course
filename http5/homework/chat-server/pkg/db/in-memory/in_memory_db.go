package in_memory

import (
	"errors"
	"sync"
)

var (
	NotExistedRowErr   = errors.New("no such row")
	NotExistedTableErr = errors.New("no such table")
)

type Table = map[string]any

type InMemDB struct {
	Tables map[string]Table
	m      *sync.RWMutex
}

func NewInMemDB() *InMemDB {
	return &InMemDB{
		Tables: make(map[string]Table),
		m:      &sync.RWMutex{},
	}
}

func (db *InMemDB) CreateTable(name string) {
	db.m.Lock()
	defer db.m.Unlock()

	db.Tables[name] = make(Table)
}

func (db *InMemDB) GetTable(name string) (Table, error) {
	db.m.RLock()
	defer db.m.RUnlock()

	t, ok := db.Tables[name]
	if ok {
		return t, nil
	}

	return nil, NotExistedTableErr
}

func (db *InMemDB) DropTable(name string) {
	db.m.Lock()
	defer db.m.Unlock()

	delete(db.Tables, name)
}

func (db *InMemDB) AlterTable(name string, newName string) {
	// needed?
}

func (db *InMemDB) Clear() {
	db.m.Lock()
	defer db.m.Unlock()

	db.Tables = make(map[string]Table)
}

func (db *InMemDB) AddRow(table string, identifier string, row any) error {
	db.m.Lock()
	defer db.m.Unlock()

	t, err := db.GetTable(table)
	if err != nil {
		return err
	}

	t[identifier] = row

	return nil
}

func (db *InMemDB) AlterRow(table string, identifier string, newRow any) error {
	db.m.Lock()
	defer db.m.Unlock()

	t, err := db.GetTable(table)
	if err != nil {
		return err
	}

	_, existed := t[identifier]
	if !existed {
		return NotExistedRowErr
	}

	t[identifier] = newRow

	return nil
}

func (db *InMemDB) GetRow(table string, identifier string) (any, error) {
	t, err := db.GetTable(table)
	if err != nil {
		return 0, err
	}

	db.m.RLock()
	defer db.m.RUnlock()

	row, exist := t[identifier]
	if !exist {
		return nil, NotExistedRowErr
	}

	return row, nil
}

func (db *InMemDB) GetAllRows(table string) ([]any, error) {
	t, err := db.GetTable(table)
	if err != nil {
		return nil, err
	}

	res := make([]any, 0, len(t))

	db.m.RLock()
	defer db.m.RUnlock()

	for _, row := range t {
		res = append(res, row)
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

	return len(t), nil
}

func (db *InMemDB) DropRow(table string, identifier string) error {
	db.m.Lock()
	defer db.m.Unlock()

	t, err := db.GetTable(table)
	if err != nil {
		return err
	}

	delete(t, identifier)

	return nil
}
