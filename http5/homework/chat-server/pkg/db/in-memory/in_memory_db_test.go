package in_memory

import (
	"testing"
)

var inMemDB = InMemDB{}

func TestTableCreated(t *testing.T) {
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

}
