package main

import (
	"fmt"
	"testing"

	_ "github.com/cznic/ql/driver"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Example struct {
	ID   uint `gorm:"primary_key"`
	Name string
	Data map[string]interface{} `gorm:"type:json"`
}

func TestJsonColumn(t *testing.T) {
	dialects := []string{"duckdb", "postgres", "sqlite3"}

	// Open a connection to each database
	for _, dialect := range dialects {
		var db *gorm.DB
		defer db.Close()
		db, err := gorm.Open(conn)
		if err != nil {
			db.Close()
			return nil, eris.Wrap(err, "failed to connect to database")
		}
		if err != nil {
			fmt.Printf("Failed to connect to %s database: %v\n", dialect, err)
			continue
		}
		fmt.Printf("Connected to %s database successfully!\n", dialect)

		// Create a simple table for the connected database
		err = createTable(db)
		if err != nil {
			fmt.Println("Failed to create the table:", err)
			return
		}
		fmt.Println("Table created successfully!")
	}
}

func openConnection(dialect string) (*gorm.DB, error) {
	// Configure the database connection based on the selected dialect
	switch dialect {
	case "postgres":
		return gorm.Open("postgres", "user=username dbname=mydb sslmode=disable")
	case "sqlite3":
		return gorm.Open("sqlite3", "test.db")
	case "duckdb":
		return gorm.Open("ql", "file:test.ql?mode=memory")
	default:
		return nil, fmt.Errorf("unsupported dialect: %s", dialect)
	}
}

func createTable(db *gorm.DB) error {
	// AutoMigrate creates the table based on the struct definition
	return db.AutoMigrate(&Example{}).Error
}
