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
		startm = startm + fmt.Sprintf("\tdb.AutoMigrate(&%s_Product{})\n", strings.Title(o.Code))
		startm = startm + fmt.Sprintf("\tdb.AutoMigrate(&%s_Product_Attributes{})\n", strings.Title(o.Code))
		startm = startm + fmt.Sprintf("\tdb.AutoMigrate(&%s_Term{})\n", strings.Title(o.Code))
		startm = startm + fmt.Sprintf("\tdb.AutoMigrate(&%s_Term_Attributes{})\n", strings.Title(o.Code))
		startm = startm + fmt.Sprintf("\tdb.AutoMigrate(&%s_Term_PriceDimensions{})\n", strings.Title(o.Code))
		startm = startm + fmt.Sprintf("\tdb.AutoMigrate(&%s_Term_PricePerUnit{})\n", strings.Title(o.Code))

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

type raw%s struct {
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]%s_Product
	Terms		map[string]map[string]map[string]raw%s_Term
}


type raw%s_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]%s_Term_PriceDimensions
	TermAttributes map[string]string
}

func (l *%s) UnmarshalJSON(data []byte) error {
	var p raw%s
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	products := []*%s_Product{}
	terms := []*%s_Term{}

	// Convert from map to slice
	for _, pr := range p.Products {
		products = append(products, &pr)
	}

	for _, tenancy := range p.Terms {
		// OnDemand, etc.
		for _, sku := range tenancy {
			// Some junk SKU
			for _, term := range sku {
				pDimensions := []*%s_Term_PriceDimensions{}
				tAttributes := []*%s_Term_Attributes{}

				for _, pd := range term.PriceDimensions {
					pDimensions = append(pDimensions, &pd)
				}

				for key, value := range term.TermAttributes {
					tr := %s_Term_Attributes{
						Key: key,
						Value: value,
					}
					tAttributes = append(tAttributes, &tr)
				}

				t := %s_Term{
					OfferTermCode: term.OfferTermCode,
					Sku: term.Sku,
					EffectiveDate: term.EffectiveDate,
					TermAttributes: tAttributes,
					PriceDimensions: pDimensions,
				}

				terms = append(terms, &t)
			}
		}
	}

	l.FormatVersion = p.FormatVersion
	l.Disclaimer = p.Disclaimer
	l.OfferCode = p.OfferCode
	l.Version = p.Version
	l.PublicationDate = p.PublicationDate
	l.Products = products
	l.Terms = terms
	return nil
}

type %s struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	[]*%s_Product `+"`gorm:\"ForeignKey:%sID\"`"+`
	Terms		[]*%s_Term`+"`gorm:\"ForeignKey:%sID\"`",
		tName, tName, tName, tName, tName,
		tName, tName, tName, tName, tName,
		tName, tName, tName, tName, tName,
		tName, tName, tName)
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

		finish = finish + o.NewStruct(fmt.Sprintf("%s_%s", tName, "Product"), p.(map[string]interface{}), tName)
	}

	// Straight up special case Terms
	finish = finish + fmt.Sprintf(`
type %s_Term struct {
	gorm.Model
	OfferTermCode string
	%sID	uint
	Sku	string
	EffectiveDate string
	PriceDimensions []*%s_Term_PriceDimensions `+"`gorm:\"ForeignKey:%s_TermID\"`"+`
	TermAttributes []*%s_Term_Attributes `+"`gorm:\"ForeignKey:%s_TermID\"`"+`
}

type %s_Term_Attributes struct {
	gorm.Model
	%s_TermID	uint
	Key	string
	Value	string
}

type %s_Term_PriceDimensions struct {
	gorm.Model
	%s_TermID	uint
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	*%s_Term_PricePerUnit `+"`gorm:\"ForeignKey:%s_Term_PriceDimensionsID\"`"+`
	// AppliesTo	[]string
}

type %s_Term_PricePerUnit struct {
	gorm.Model
	%s_Term_PriceDimensionsID	uint
	USD	string
}`, tName, tName, tName, tName, tName, tName, tName, tName, tName, tName, tName, tName, tName, tName)

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

func (o Structure) NewStruct(name string, val map[string]interface{}, parent string) string {
	tName := strings.Title(name)
	start := fmt.Sprintf(`type %s struct {
	gorm.Model
	`, tName)
	middling := fmt.Sprintf("\t%s\tuint\n", parent+"ID")
	finish := "}\n"

	for field, value := range val {
		fUpper := strings.Title(field)
		entry := ""
		switch value.(type) {
		case map[string]interface{}:
			n := fmt.Sprintf("%s_%s", tName, fUpper)
			entry = fmt.Sprintf("\t%s\t%s\t`gorm:\"ForeignKey:%s\"`", fUpper, n, n+"ID")
			isGarbage := false
			for _, r := range fUpper {
				switch {
				case r >= '0' && r <= '9':
					isGarbage = true
				}
			}
			if !isGarbage {
				finish = finish + o.NewStruct(n, value.(map[string]interface{}), n)
			} else {
				counter := 0
				for _, k := range value.(map[string]interface{}) {
					n = n + "\t`gorm:\"ForeignKey:ID\"`"
					if counter > 0 {
						break
					}
					counter++
					v, ok := k.(map[string]interface{})
					if !ok {
						entry = fmt.Sprintf("\t%s\t%s", entry, fUpper, "string")
					} else {
						finish = finish + o.NewStruct(n, v, n)
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
