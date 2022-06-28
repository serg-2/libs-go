package databaselib

import (
	"database/sql"
	"fmt"
	cl "github.com/serg-2/libs-go/commonlib"
	"log"
	"time"
)

// Database - typical database structure
type Database struct {
	Instance *sql.DB
	Rows     *sql.Rows
}

// NewDb - initialize new database
func (d *Database) NewDb(config cl.DatabaseConfig, applicationName string) {

	// Generate credentials
	credentials := fmt.Sprintf("user=%s password=%s dbname=%s host=%s application_name=%s", config.Username, config.Password, config.Database, config.Host, applicationName)

	//Open DATABASE
	var err error
	log.Println("Connecting to database...")
	d.Instance, err = sql.Open("postgres", credentials)
	cl.ChkFatal(err)
	d.Instance.SetMaxOpenConns(70)
	log.Println("Connected to database!")
}

// DbClose - Close database
func (d *Database) DbClose() {
	d.Instance.Close()
}

// Check - Check Database
func (d *Database) Check() {
	log.Println("Checking DB...")
	err := d.Instance.Ping()
	if err != nil {
		log.Println("Checking DB - FAIL! Error:")
		log.Println(err)
		for {
			log.Println("Reconnecting to DB...")
			err2 := d.Instance.Ping()
			if err2 != nil {
				time.Sleep(5 * time.Second)
			} else {
				log.Println("Reconnect to DB successfully.")
				break
			}
		}
	}
	log.Println("Checking DB - SUCCESS!")
}

// GetTx - Get TX from database
func (d *Database) GetTx() *sql.Tx {
	tx, err := d.Instance.Begin()
	cl.ChkFatal(err)
	return tx
}

// CheckTableExists - check table exists
func (d *Database) CheckTableExists(table_name string) bool {
	// Check table exists
	query := "SELECT count(table_name) FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name = '" + table_name + "'"
	// log.Println("Querying database with query from Check Table Exists: " + query)

	var countRows int

	err := d.Instance.QueryRow(query).Scan(&countRows)
	cl.ChkFatal(err)

	if countRows != 1 {
		return false
	}

	return true
}
