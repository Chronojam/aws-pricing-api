package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var baseUrl = "https://pricing.us-east-1.amazonaws.com"

type offer struct {
	Code              string `json:"offerCode"`
	CurrentVersionUrl string `json:"currentVersionUrl"`
}
type offerResponse struct {
	Offers map[string]offer `json:"offers"`
}

func main() {
	resp, err := http.Get(baseUrl + "/offers/v1.0/aws/index.json")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var or offerResponse
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(b, &or)

	// Create a global functions file
	startm := fmt.Sprintf(`
package schema

import (
	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB) error {
`)

	finishm := `
	return nil
}`

	for _, o := range or.Offers {
		fmt.Println("Writing: " + o.Code)
		res, err := http.Get(baseUrl + o.CurrentVersionUrl)
		if err != nil {
			panic(err)
		}
		b, err = ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		raw := map[string]interface{}{}
		err = json.Unmarshal(b, &raw)
		if err != nil {
			panic(err)
		}

		startm = startm + fmt.Sprintf("\tdb.AutoMigrate(&%s{})\n", strings.Title(o.Code))
		ProcessForSchema(raw, o.Code, baseUrl+o.CurrentVersionUrl)
		// ioutil.WriteFile("./raw/"+o.Code+".json", b, 0655)
	}

	finalm := startm + finishm
	ioutil.WriteFile("./schema/global.go", []byte(finalm), 0655)
}

type Structure map[string]interface{}

func (o Structure) Test(name string, val map[string]interface{}, url string) string {
	tName := strings.Title(name)
	start := fmt.Sprintf(`package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type %s struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]%s_Product
	Terms		map[string]map[string]map[string]%s_Term`, tName, tName, tName)
	middling := "\n"
	finish := "}\n"

	// We know that Products and Terms will always exist.
	counter := 0
	for _, p := range val["products"].(map[string]interface{}) {
		// We just want the product schema, so take the first one.
		if counter > 0 {
			break
		}
		counter++

		finish = finish + o.NewStruct(fmt.Sprintf("%s_%s", tName, "Product"), p.(map[string]interface{}))
	}

	// Straight up special case Terms
	finish = finish + fmt.Sprintf(`
type %s_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]%s_Term_PriceDimensions
	TermAttributes map[string]string
}

type %s_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	%s_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type %s_Term_PricePerUnit struct {
	USD	string
}`, tName, tName, tName, tName, tName, tName)

	// Add some helper functions to pull api data.
	finish = finish + fmt.Sprintf(`
func (a %s) QueryProducts(q func(product %s_Product) bool) []%s_Product{
	ret := []%s_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}`, tName, tName, tName, tName)

	finish = finish + fmt.Sprintf(`
func (a %s) QueryTerms(t string, q func(product %s_Term) bool) []%s_Term{
	ret := []%s_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}`, tName, tName, tName, tName)

	finish = finish + fmt.Sprintf(`
func (a *%s) Refresh() error {
	var url = "%s"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, a)
	if err != nil {
		return err
	}

	return nil
}`, tName, url)

	return start + middling + finish
}

func (o Structure) NewStruct(name string, val map[string]interface{}) string {
	tName := strings.Title(name)
	start := fmt.Sprintf("type %s struct {", tName)
	middling := ""
	finish := "}\n"

	for field, value := range val {
		fUpper := strings.Title(field)
		entry := ""
		switch value.(type) {
		case map[string]interface{}:
			n := fmt.Sprintf("%s_%s", tName, fUpper)
			entry = fmt.Sprintf("\t%s\t%s", fUpper, n)
			isGarbage := false
			for _, r := range fUpper {
				switch {
				case r >= '0' && r <= '9':
					isGarbage = true
				}
			}
			if !isGarbage {
				finish = finish + o.NewStruct(n, value.(map[string]interface{}))
			} else {
				counter := 0
				for _, k := range value.(map[string]interface{}) {
					if counter > 0 {
						break
					}
					counter++
					v, ok := k.(map[string]interface{})
					if !ok {
						entry = fmt.Sprintf("\t%s\t%s", fUpper, "string")
					} else {
						finish = finish + o.NewStruct(n, v)
					}
				}
			}
		case string:
			entry = fmt.Sprintf("\t%s\t%s", fUpper, "string")
		}

		middling = middling + entry + "\n"
	}
	return start + middling + finish
}

func ProcessForSchema(raw map[string]interface{}, code string, url string) {
	obj := Structure{}
	out := obj.Test(code, raw, url)
	ioutil.WriteFile("./schema/"+strings.Title(code)+".go", []byte(out), 0655)
}
