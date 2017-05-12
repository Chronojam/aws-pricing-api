package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/chronojam/aws-pricing-api/types/schema"
)

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = schema.Migrate(db)
	if err != nil {
		panic(err)
	}
}
