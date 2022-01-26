// Copyright Kueski. All rights reserved.
// Use of this source code is not licensed

// Package db provides database connection services 
package db

import (
	"context"
	"fmt"
	"log"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/gastonstec/utils/config"
)

// Connection pool
var DBpool	*pgxpool.Pool

// Opens a connection pool
func OpenDB() error {
	
	// create connection
	var err error = nil
	DBpool, err = pgxpool.Connect(context.Background(), config.GetString("DBURI"))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// check connection and database info
	rows, err := DBpool.Query(context.Background(), "select version()::TEXT as version, current_database()::TEXT as database, inet_server_addr()::TEXT as server, inet_server_port()::TEXT as port")

	// check for errors
	if (err != nil) {
    	log.Println(err.Error())
		return err
	} else {		
		// get database info
		var version, database, server, port string
		rows.Next()
		rows.Scan(&version, &database, &server, &port)
		rows.Close()
		// log database info
		log.Println(fmt.Sprintf("Connected to dbname=%s version=%s on server=%s port=%s",
						database, version, server, port))
	}

	return nil
}

// Closes the connection pool
func CloseDB() {
	DBpool.Close()
	log.Println("Database closed")
}