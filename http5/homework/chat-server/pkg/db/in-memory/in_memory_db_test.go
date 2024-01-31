package in_memory

import (
	"context"
	"testing"
)

var ctx, cancel = context.WithCancel(context.Background())
var inMemDB = NewInMemDB(ctx, "db_save.json")

func TestTableCreated(t *testing.T) { // todo: make it pass
	inMemDB.Clear()

	tableName := "new_table"

	inMemDB.CreateTable(tableName)

	if _, ok := inMemDB.Tables[tableName]; !ok {
		t.Fatal()
	}
}

func TestGetExistingTable(t *testing.T) {
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
