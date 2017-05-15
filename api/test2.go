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
