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

var imDB *memdb.MemDB

type Network struct {
	Code	string
	Name 	string
}

func CreateSchema() (*memdb.DBSchema) {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
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

func LoadNetwork(imDB *memdb.MemDB) (int64, int64, error) {
	var dbRecords, memRecords int64 = 0, 0
	var err error

	row := db.DBpool.QueryRow(context.Background(), "SELECT count(network_code) FROM network")
	err = row.Scan(&dbRecords)
	if dbRecords <= 0 {
		return dbRecords, 0, err
	}


	rows, err := db.DBpool.Query(context.Background(), "SELECT network_code, network_name FROM network ORDER BY network_code")
	if err != nil {
		return dbRecords, memRecords, err
	}
	defer rows.Close()


	txn := imDB.Txn(true)

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


	// Commit the transaction
	txn.Commit()

	return dbRecords, memRecords, nil
}

func Load() error {
	var err error

	// Create schema
	schema := CreateSchema()
	imDB, err = memdb.NewMemDB(schema)
	if err != nil {
		return err
	}

	// Load network table 
	var dbRecords, memRecords int64 = 0, 0
	dbRecords, memRecords, err = LoadNetwork(imDB)
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


func GetNetwork (categoryCode string) (string, error) {

	txn := imDB.Txn(false)
	defer txn.Abort()
	
	raw, err := txn.First("network", "id", categoryCode)
	if err != nil {
		return "", err
	}

	return raw.(Network).Name, nil

}