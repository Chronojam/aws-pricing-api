package main

import (
	"fmt"
	"github.com/chronojam/aws-pricing-api/types/schema"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db, err := gorm.Open("mysql", "root:my-secret-pw@/pricing?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	//	db.Exec("PRAGMA foreign_keys = ON")
	//	db.LogMode(true)

	defer db.Close()

	err = schema.Migrate(db)
	if err != nil {
		panic(err)
	}

	/*	p := &schema.AWSBudgets{}
		err = p.Refresh()
		if err != nil {
			panic(err)
		}
		db.Create(p)
	*/

	l := &schema.AWSBudgets{}

	db.First(l)

	db.Model(l).Association("Products").Find(&l.Products)
	//db.Model(l).Association("Terms").Find(l.Terms)
	for _, t := range l.Products {
		fmt.Printf("%v:\n", t)
		db.First(t)
		db.Model(t).Association("Attributes").Find(&t.Attributes)
		fmt.Printf("\t%v\n", t.Attributes)
	}

	db.Model(l).Association("Terms").Find(&l.Terms)
	//db.Model(l).Association("Terms").Find(l.Terms)
	for _, t := range l.Terms {
		fmt.Printf("%v:\n", t)
		db.First(t)
		db.Model(t).Association("PriceDimensions").Find(&t.PriceDimensions)
		for _, a := range t.PriceDimensions {
			fmt.Printf("\t%v\n", a)
		}
	}
}
