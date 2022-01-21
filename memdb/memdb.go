// Package models implements an in-memory database
package memdb

import (
	"fmt"
	"errors"
	"log"
	"context"
	"strings"
	"github.com/gastonstec/utils"
	"gastonstec/nuricc/db"
	"github.com/hashicorp/go-memdb"
)

// In-memory database pointer
var imDB *memdb.MemDB

// Blockchain network struct
type Network struct {
	Code	string
	Name 	string
}

// Function CreateSchema creates the memory database schema
func CreateSchema() (*memdb.DBSchema) {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			// Network table
			"network": {
				Name: "network",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Code"},
					},
					"name": {
						Name:    "name",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Name"},
					},
				},
			},
		},
	}

	return schema
}

// Function LoadNetwork loads the memory database schema with
// values stored in the database
func loadNetwork(imDB *memdb.MemDB) (int64, int64, error) {
	var dbRecords, memRecords int64 = 0, 0
	var err error

	// Get the total number of records in the database
	row := db.DBpool.QueryRow(context.Background(), "SELECT count(network_code) FROM network")
	err = row.Scan(&dbRecords)
	if dbRecords <= 0 {
		return dbRecords, 0, err
	}

	// Get the records from the database
	rows, err := db.DBpool.Query(context.Background(), 
					"SELECT network_code, network_name FROM network ORDER BY network_code")
	if err != nil {
		return dbRecords, memRecords, err
	}
	defer rows.Close()

	// Insert records on the memory database
	txn := imDB.Txn(true) // write transaction
	var nw Network
	for rows.Next() {
		
		err = rows.Scan(&nw.Code, &nw.Name)
		if err != nil {
			return dbRecords, memRecords, err
		}

		nw.Code = strings.TrimSpace(nw.Code)
		nw.Name = strings.TrimSpace(nw.Name)

		err = txn.Insert("network", nw)
		if err != nil {
			return dbRecords, memRecords, err
		}

		memRecords += 1
	}
	txn.Commit() // Commit the transaction

	// Return control values
	return dbRecords, memRecords, nil
}

// Function Load loads the memory database 
func Load() error {
	var err error

	// Create memory database schema
	schema := CreateSchema()
	imDB, err = memdb.NewMemDB(schema)
	if err != nil {
		return err
	}

	// Load network table from the database
	var dbRecords, memRecords int64 = 0, 0
	dbRecords, memRecords, err = loadNetwork(imDB)
	if err != nil {
		return err
	}
	// Check loaded records
	if dbRecords > memRecords {
		return errors.New("network table loaded with fewer database records")
	}
	// Table loaded ok
	log.Println(fmt.Sprintf("%s: network table loaded with %d records of %d", utils.GetFunctionName(), memRecords, dbRecords))

	return nil
}

// Function GetNetwork returns the network name using
// the network code
func GetNetwork (categoryCode string) (string, error) {

	txn := imDB.Txn(false) // read transaction
	defer txn.Abort()
	raw, err := txn.First("network", "id", categoryCode)
	if err != nil || raw == nil {
		return "", err
	}

	return raw.(Network).Name, nil
}