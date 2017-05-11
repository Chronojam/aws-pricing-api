package main

import (
	"fmt"
	"github.com/chronojam/aws-pricing-api/types/schema"
	"strconv"
	"strings"
	"testing"
)

type Thing struct {
	Name  string
	VCpu  string
	Mem   string
	Price string
}

func TestLel(t *testing.T) {
	p := &schema.AmazonEC2{}
	err := p.Refresh()
	if err != nil {
		fmt.Println(err.Error())
	}

	// Find our product
	//	vals := map[string]schema.AmazonEC2_Product{}
	l := p.QueryProducts(func(prod schema.AmazonEC2_Product) bool {
		return prod.Attributes.Tenancy == "Shared" && // Dedicated, Host
			prod.Attributes.Location == "EU (Ireland)" &&
			prod.Attributes.OperatingSystem == "Linux"
	})

	tList := []*Thing{}

	for _, prod := range l {
		t := &Thing{
			VCpu: prod.Attributes.Vcpu,
			Mem:  prod.Attributes.Memory,
			Name: prod.Attributes.InstanceType,
		}

		onDemandTerms := p.QueryTerms("OnDemand", func(term schema.AmazonEC2_Term) bool {
			return term.Sku == prod.Sku
		})

		for _, v := range onDemandTerms {
			for _, price := range v.PriceDimensions {
				t.Price = price.PricePerUnit.USD
			}
		}

		tList = append(tList, t)
	}

	cheapestCpu := 1.0
	cObj := &Thing{}
	cheapestMem := 1.0
	mObj := &Thing{}

	for _, i := range tList {
		cpu, err := strconv.ParseFloat(i.VCpu, 64)
		if err != nil {
			panic(err)
		}
		price, err := strconv.ParseFloat(i.Price, 64)
		if err != nil {
			panic(err)
		}
		pm := strings.Replace(strings.Replace(strings.Replace(i.Mem, ",", "", -1), " Gib", "", -1), " GiB", "", -1)
		mem, err := strconv.ParseFloat(pm, 64)
		if err != nil {
			panic(err)
		}

		cpuValue := price / cpu
		memValue := price / mem
		if cpuValue < cheapestCpu {
			cheapestCpu = cpuValue
			cObj = i
		}

		if memValue < cheapestMem {
			cheapestMem = memValue
			mObj = i
		}

		fmt.Printf(i.Name + "\n")
	}

	fmt.Printf("Cheapest CPU: %v\n", cObj.Name)
	fmt.Printf("@Value: %v\n", cheapestCpu)
	fmt.Printf("Cheapest Mem: %v\n", mObj.Name)
	fmt.Printf("@Value: %v\n", cheapestMem)
}
