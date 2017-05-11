package schema

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/jinzhu/gorm"
)

type AmazonLex struct {
	gorm.Model
	FormatVersion	string
	Disclaimer	string
	OfferCode	string
	Version		string
	PublicationDate	string
	Products	map[string]AmazonLex_Product
	Terms		map[string]map[string]map[string]AmazonLex_Term
}
type AmazonLex_Product struct {	Sku	string
	ProductFamily	string
	Attributes	AmazonLex_Product_Attributes
}
type AmazonLex_Product_Attributes struct {	OutputMode	string
	Location	string
	LocationType	string
	Group	string
	Operation	string
	SupportedModes	string
	Servicecode	string
	GroupDescription	string
	Usagetype	string
	InputMode	string
}

type AmazonLex_Term struct {
	OfferTermCode string
	Sku	string
	EffectiveDate string
	PriceDimensions map[string]AmazonLex_Term_PriceDimensions
	TermAttributes AmazonLex_Term_TermAttributes
}

type AmazonLex_Term_PriceDimensions struct {
	RateCode	string
	RateType	string
	Description	string
	BeginRange	string
	EndRange	string
	Unit	string
	PricePerUnit	AmazonLex_Term_PricePerUnit
	AppliesTo	[]interface{}
}

type AmazonLex_Term_PricePerUnit struct {
	USD	string
}

type AmazonLex_Term_TermAttributes struct {

}
func (a AmazonLex) QueryProducts(q func(product AmazonLex_Product) bool) []AmazonLex_Product{
	ret := []AmazonLex_Product{}
	for _, v := range a.Products {
		if q(v) {
			ret = append(ret, v)
		}
	}

	return ret
}
func (a AmazonLex) QueryTerms(t string, q func(product AmazonLex_Term) bool) []AmazonLex_Term{
	ret := []AmazonLex_Term{}
	for _, v := range a.Terms[t] {
		for _, val := range v {
			if q(val) {
				ret = append(ret, val)
			}
		}
	}

	return ret
}
func (a *AmazonLex) Refresh() error {
	var url = "https://pricing.us-east-1.amazonaws.com/offers/v1.0/aws/AmazonLex/current/index.json"
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
}