package main

import (
	"fmt"
	"github.com/chronojam/aws-pricing-api/types/schema"
	"testing"
)

func TestLel(t *testing.T) {
	p := &schema.AmazonEC2{}
	err := p.Refresh()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(p.FormatVersion)
}
