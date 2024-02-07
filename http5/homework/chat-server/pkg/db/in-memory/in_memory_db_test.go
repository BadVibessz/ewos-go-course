package in_memory

import (
	"context"
	"testing"
)

func initDB() *InMemDB {
	ctx := context.Background()
	inMemDB, _ := NewInMemDB(ctx, "")

	return inMemDB
}

func TestTableCreated(t *testing.T) {
	inMemDB := initDB()

	inMemDB.Clear()

	tableName := "new_table"

	inMemDB.CreateTable(tableName)

	if _, ok := inMemDB.Tables[tableName]; !ok {
		t.Fatal()
	}
}

func TestGetExistingTable(t *testing.T) {
	inMemDB := initDB()

	inMemDB.Clear()

	tableName := "new_table"

	inMemDB.CreateTable(tableName)

	_, err := inMemDB.GetTable(tableName)
	if err != nil {
		t.Fatal()
	}

	tableName = "new_table2"

	inMemDB.CreateTable(tableName)

	_, err = inMemDB.GetTable(tableName)
	if err != nil {
		t.Fatal()
	}
}
