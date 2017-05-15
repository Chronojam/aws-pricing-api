AWS-Pricing-API
--

This library allows you to query the amazon pricing api using go types
It can cache the results in a database (using gorm)


In memory usage:
=
``` go
package main

import (
        "fmt"
        "github.com/chronojam/aws-pricing-api/types/schema"
        "strings"
)

func main() {
        ec2 := &schema.AmazonEC2{}

        // Populate this object with new pricing data
        err := ec2.Refresh()
        if err != nil {
                panic(err)
        }

        // Get the price of all c4.Large instances,
        // running linux, on shared tenancy
        c4Large := []*schema.AmazonEC2_Product{}
        for _, p := range ec2.Products {
                if p.Attributes.InstanceType == "c4.large" &&
                        p.Attributes.OperatingSystem == "Linux" &&
                        p.Attributes.Tenancy == "Shared" {
                        c4Large = append(c4Large, p)
                }
        }

        // Show the pricing data for each of those.
        for _, p := range c4Large {
                //fmt.Println(p.Sku)
                // Find the correct terms
                for _, term := range ec2.Terms {
                        if term.Sku == p.Sku {
                                for _, pd := range term.PriceDimensions {
                                        // I Stripped out the OnDemand/Reserved field, but maybe ill add it back later
                                        // Only On Demand
                                        if strings.Contains(pd.Description, "On Demand") {
                                                fmt.Printf("%s:\n", p.Sku)
                                                fmt.Printf("\t%s:\n", "PriceDimensions")
                                                fmt.Printf("\t\t%s\n", pd.Description)
                                                fmt.Printf("\t\t%s\n", pd.PricePerUnit.USD)
                                        }
                                }
                        }
                }
        }
}
```

With a backing DB
=
``` go

import (
        "fmt"
        "github.com/chronojam/aws-pricing-api/types/schema"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)


// Call me to populate with the appropriate product
// Dont call me twice for the same product, because you'll end up with duplicate entries.
func Populate(db *gorm.DB) {
     // p := &schema.AmazonEC2{} or whatever you need
      p := &schema.AWSBudgets{}
      err = p.Refresh()
      if err != nil {
	      panic(err)
      }
      db.Create(p)
}

func main() {
        db, err := gorm.Open("mysql", "root:my-secret-pw@/pricing?charset=utf8&parseTime=True&loc=Local")
        if err != nil {
                panic(err)
        }

        defer db.Close()

        err = schema.Migrate(db)
        if err != nil {
                panic(err)
        }

	// Populate(db)

	// Do some querying
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
```


Updating generated types
=

This'll take a little while, as it pulls down everything
```bash
cd types/
make generate
```
